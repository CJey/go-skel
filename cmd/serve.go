package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cjey/overseer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/cjey/gbase"

	"go-skel/app"
)

const (
	LISTEN_PORT      = 1234 // --listen default port
	GRPC_LISTEN_PORT = 4321 // --listen-grpc default port
)

var (
	_ = fmt.Print

	_Counter = struct {
		Tcp  int64
		Http int64
		Grpc int64
	}{}

	_CMDServe = &cobra.Command{
		Use:   `serve`,
		Run:   runServe,
		Short: `Running a demo http service`,
		Long:  `Running a demo http service`,
	}
)

func init() {
	var cmd = _CMDServe
	supportConfigAndLogger(cmd)

	// --listen 0.0.0.0:LISTEN_PORT
	cmd.PersistentFlags().String("listen", strconv.Itoa(LISTEN_PORT),
		"listening address for demo Http service")
	// --listen-grpc 0.0.0.0:GRPC_LISTEN_PORT
	cmd.PersistentFlags().String("listen-grpc", strconv.Itoa(GRPC_LISTEN_PORT),
		"listening address for demo Grpc service")

	// do not bind, flag only
	// --nake=false
	cmd.PersistentFlags().Bool("nake", false,
		"do not use overseer to protect service")
	// --grace-signal SIGUSR2
	cmd.PersistentFlags().String("grace-signal", "SIGUSR2",
		"which signal used to trigger reload/shutdown gracefully, only from flag, support SIGUSR1, SIGUSR2, SIGINT")

	_CMDRoot.AddCommand(cmd)
}

func runServe(cmd *cobra.Command, args []string) {
	handleConfigAndLogger(cmd)

	// bind flags
	viper.BindPFlag("listen", cmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("listen-grpc", cmd.PersistentFlags().Lookup("listen-grpc"))

	var naked, _ = cmd.Flags().GetBool("nake")
	if naked == false && overseer.IsSupported() && runtime.GOOS == "linux" {
		runServeWithOverseer(cmd, args)
	} else {
		runServeDirectly(cmd, args)
	}
}

func runServeWithOverseer(cmd *cobra.Command, args []string) {
	var ctx gbase.Context
	// WARNING: unexported method at overseer source code
	var isMaster = os.Getenv("OVERSEER_IS_SLAVE") != "1"

	if isMaster {
		ctx = gbase.NamedContext("overseer-master$" + strconv.Itoa(os.Getpid()))
	} else {
		ctx = gbase.NamedContext("overseer-worker#" + os.Getenv("OVERSEER_SLAVE_ID") + "$" + strconv.Itoa(os.Getpid()))
	}
	defer ctx.Logger().Sync()

	var addrHTTP, addrGRPC, sigGrace = serveCheckBasicConfig(ctx, cmd)

	var ocfg = overseer.Config{
		Addresses:        []string{addrHTTP, addrGRPC},
		Required:         true,
		RestartSignal:    sigGrace,
		TerminateTimeout: 9999 * time.Hour, // do not terminate

		Program: func(state overseer.State) {
			var (
				lsnHTTP     = state.Listeners[0]
				lsnGRPC     = state.Listeners[1]
				apptitle    = app.Name + "#" + os.Getenv("OVERSEER_SLAVE_ID")
				ctx, cancel = gbase.NamedContext(
					app.Name + "#" + os.Getenv("OVERSEER_SLAVE_ID") + "$" + strconv.Itoa(os.Getpid()),
				).WithCancel()
			)
			defer ctx.Logger().Sync()

			// reload signal received
			go func() {
				<-state.GracefulShutdown
				ctx.Info("Graceful shutdown signal received")

				// diff * with #, means shutting down
				apptitle = app.Name + "*" + os.Getenv("OVERSEER_SLAVE_ID")
				cancel()
			}()

			// signal received by self
			go func() {
				var (
					cnt   int
					last  time.Time
					chsig = make(chan os.Signal)
				)

				signal.Notify(chsig, syscall.SIGHUP, syscall.SIGTERM)
				if sigGrace != syscall.SIGINT {
					signal.Notify(chsig, syscall.SIGINT)
				}

				for {
					var sig = <-chsig
					if sig == syscall.SIGTERM {
						var p, _ = os.FindProcess(os.Getpid())
						// must send sigGrace to myself, not my parent
						// we can not cancel ctx directly, because of the listener inherited from overseer work bad
						p.Signal(sigGrace)
					} else {
						if sig == syscall.SIGINT {
							// Ctrl+C * 2+ fastly. Yes sir, force quit!
							if time.Since(last) > 500*time.Millisecond {
								cnt = 1
							} else {
								cnt++
								if cnt >= 2 {
									ctx.Warn("Fastly Ctrl+C action detected, force exit!")
									os.Exit(1)
								}
							}
							last = time.Now()
						}
						ctx.Warn("Signal received and ignored",
							"signal", signalName(sig), "grace-signal", signalName(sigGrace), "ppid", os.Getppid())
					}
					// avoid too fast
					time.Sleep(5 * time.Millisecond)
				}
			}()

			// use process name to expose core runtime information every 1s
			gbase.LiveProcessName(ctx, 1*time.Second, func(int) string {
				return fmt.Sprintf("%s [T%d H%d G%d] v%s", apptitle,
					_Counter.Tcp, _Counter.Http, _Counter.Grpc,
					app.FullVersion,
				)
			})

			// serve
			serve(ctx, cmd, args, lsnHTTP, lsnGRPC)
		},
	}

	if isMaster {
		gbase.WriteMyPID(ctx, app.Name)
		// sigGrace send to overseer master, trigger gracefully restart
		// sigGrace send to overseer worker, trigger gracefully shutdown
		ctx.Info("Use signal " + signalName(sigGrace) + " to reload service gracefully")

		// TODO: other master care things, e.g. zk register
	}

	// listen & serve
	for {
		var err = overseer.RunErr(ocfg)
		if err == nil {
			ctx.Info("Exited")
			break
		}

		if isMaster {
			if e, ok := err.(*net.OpError); ok && e.Op == "listen" {
				ctx.Fatal("Listen failed", "addrs", ocfg.Addresses, "err", err)
				break
			}
			ctx.Error("Master failed, waiting to restart", "err", err)
		} else {
			ctx.Error("Worker start failed, waiting to retry", "err", err)
		}

		// unexpected error, must retry
		time.Sleep(3 * time.Second)
	}
}

func runServeDirectly(cmd *cobra.Command, args []string) {
	var ctx, cancel = gbase.NamedContext(app.Name + "$" + strconv.Itoa(os.Getpid())).WithCancel()
	defer ctx.Logger().Sync()
	ctx.Warn("Running at naked mode, without overseer's watching")

	var addrHTTP, addrGRPC, sigGrace = serveCheckBasicConfig(ctx, cmd)

	// I am the master always
	// wirte pid
	gbase.WriteMyPID(ctx, app.Name)
	// TODO: other master care things, e.g. zk register

	// listen
	lsnHTTP, err := net.Listen("tcp", addrHTTP)
	if err != nil {
		ctx.Fatal("Listen http service failed", "err", err, "addr", addrHTTP)
	}
	lsnGRPC, err := net.Listen("tcp", addrGRPC)
	if err != nil {
		ctx.Fatal("Listen grpc service failed", "err", err, "addr", addrGRPC)
	}

	// signal received
	go func() {
		var (
			cnt   int
			last  time.Time
			once  sync.Once
			chsig = make(chan os.Signal)
		)

		signal.Notify(chsig, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

		for {
			var sig = <-chsig
			if sig == sigGrace || sig == syscall.SIGTERM {
				once.Do(func() {
					ctx.Info("Signal received, shutting down gracefully", "signal", signalName(sig))
					cancel()
				})
			} else {
				if sig == syscall.SIGINT {
					// Ctrl+C * 2+ fastly. Yes sir, force quit!
					if time.Since(last) > 500*time.Millisecond {
						cnt = 1
					} else {
						cnt++
						if cnt >= 2 {
							ctx.Warn("Fastly Ctrl+C action detected, force exit!")
							os.Exit(1)
						}
					}
					last = time.Now()
				}
				ctx.Warn("Signal received and ignored",
					"signal", signalName(sig), "grace-signal", signalName(sigGrace))
			}
			// avoid too fast
			time.Sleep(5 * time.Millisecond)
		}
	}()

	// serve
	serve(ctx, cmd, args, lsnHTTP, lsnGRPC)
}

func serveCheckBasicConfig(ctx gbase.Context, cmd *cobra.Command) (string, string, syscall.Signal) {
	// check config --listen
	addrHTTP, err := gbase.ResolveTCPAddr(ctx, viper.GetString("listen"), LISTEN_PORT)
	if err != nil {
		ctx.Fatal("Invalid address", "err", err, "listen", viper.GetString("listen"))
	}
	// check config --listen-grpc
	addrGRPC, err := gbase.ResolveTCPAddr(ctx, viper.GetString("listen-grpc"), GRPC_LISTEN_PORT)
	if err != nil {
		ctx.Fatal("Invalid address", "err", err, "listen", viper.GetString("listen-grpc"))
	}

	// check flag --grace-signal
	var sigGrace syscall.Signal
	var sig = cmd.Flag("grace-signal").Value.String()
	switch sig {
	case "USR1", "SIGUSR1":
		sigGrace = syscall.SIGUSR1
	case "USR2", "SIGUSR2":
		sigGrace = syscall.SIGUSR2
	case "INT", "SIGINT":
		sigGrace = syscall.SIGINT
	default:
		ctx.Fatal("Unsupported reload signal", "signal", sig)
	}

	return addrHTTP.String(), addrGRPC.String(), sigGrace
}

func serve(ctx gbase.Context, cmd *cobra.Command, args []string, lsnHTTP, lsnGRPC net.Listener) {
	// TODO: start all required procedures

	var (
		wg   sync.WaitGroup
		once sync.Once
		end  = make(chan struct{})
		boom = func() { once.Do(func() { close(end) }) }
	)

	var isClosedListener = func(err error) bool {
		return strings.Contains(err.Error(), "use of closed network connection")
	}

	// waiting ctx
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer boom()

		select {
		case <-end:
		case <-ctx.Done():
		}
	}()

	// http service [1/2]
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer boom()

		// create
		var ctx = gbase.NamedContext(ctx.Name() + ".http")
		var srv = http.Server{
			IdleTimeout:       60 * time.Second,
			WriteTimeout:      30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			ConnState: func(conn net.Conn, state http.ConnState) {
				switch state {
				case http.StateNew:
					atomic.AddInt64(&_Counter.Tcp, 1)
				case http.StateClosed:
					atomic.AddInt64(&_Counter.Tcp, -1)
				}
			},
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt64(&_Counter.Http, 1)
				defer atomic.AddInt64(&_Counter.Http, -1)
				var ctx = gbase.SessionContext()

				ctx.Info("new HTTP request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
				defer ctx.Info("end HTTP request")

				w.Write([]byte(fmt.Sprintf("%s Demo", time.Now().String())))
			}),
		}

		// gracefully shutdown
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-end
			ctx.Info("Shutting down", "listen", lsnHTTP.Addr().String())
			// listener will be closed immediately
			// idle connection will be closed immediately
			srv.Shutdown(ctx)
		}()

		// serve
	loop:
		for {
			ctx.Info("Running", "listen", lsnHTTP.Addr().String())
			if err := srv.Serve(lsnHTTP); err != nil {
				time.Sleep(50 * time.Millisecond) // wait end unblock
				select {
				case <-end:
				default:
					if !isClosedListener(err) {
						ctx.Error("HTTP serve failed, waiting to restart", "err", err)
						time.Sleep(5 * time.Second)
						continue loop
					}
				}
			}
			ctx.Info("Serve done", "listen", lsnHTTP.Addr().String())
			break
		}
	}()

	// grpc service [2/2]
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer boom()

		var interceptor = func(gctx context.Context, req interface{},
			info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rep interface{}, err error) {
			atomic.AddInt64(&_Counter.Grpc, 1)
			defer atomic.AddInt64(&_Counter.Grpc, -1)

			return handler(gctx, req)
		}

		// create
		var ctx = gbase.NamedContext(ctx.Name() + ".grpc")
		var srv = grpc.NewServer(
			grpc.UnaryInterceptor(interceptor),
			grpc.ConnectionTimeout(5*time.Second),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: 60 * time.Second,
				MaxConnectionAge:  300 * time.Second,
			}),
		)

		// gracefully shutdown
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-end
			ctx.Info("Shutting down", "listen", lsnGRPC.Addr().String())
			// listener will be closed immediately
			// idle connection will be closed immediately
			// active connection will send goaway after MaxConnectionAge
			srv.GracefulStop()
		}()

		// serve
	loop:
		for {
			ctx.Info("Running", "listen", lsnGRPC.Addr().String())
			if err := srv.Serve(lsnGRPC); err != nil {
				time.Sleep(50 * time.Millisecond) // wait end unblock
				select {
				case <-end:
				default:
					if !isClosedListener(err) {
						ctx.Error("GRPC serve failed, waiting to restart", "err", err)
						time.Sleep(5 * time.Second)
						continue loop
					}
				}
			}
			ctx.Info("Serve done", "listen", lsnGRPC.Addr().String())
			break
		}
	}()

	ctx.Info("go-skel serve demo service started", "version", app.FullVersion)
	wg.Wait()
	ctx.Info("go-skel serve demo service stopped", "version", app.FullVersion)

	//var bgonce sync.Once
	//for {
	//	// TODO: other background tasks that must not be interrupted
	//	var bgtasks = numbgtasks()
	//	if bgtasks == 0 {
	//		ctx.Info("All background tasks stopped")
	//		break
	//	}
	//	bgonce.Do(func() {
	//		ctx.Warn("Waiting for background tasks stop", "bgtasks", bgtasks)
	//	})
	//	time.Sleep(1 * time.Second)
	//}
}
