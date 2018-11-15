package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"{=APPNAME=}/build"
)

// base package, always hold on to anti golang build warning
func init() { return; fmt.Print(); zap.S() }

func main() {
	initDefaultLogger()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   build.Appname(),
	Short: ``,
	Long: ``,
	Version:          build.Version(),
	PersistentPreRun: PersistentPreRun,
	Run:              Run,
	//TraverseChildren: true,
}

func init() {
	rootCmd.PersistentFlags().StringP(
		"config", "c", build.Appname()+".toml",
		"configuration file path",
	)
	rootCmd.PersistentFlags().StringP(
		"log-file", "l", "stderr",
		"log file path, support stderr, stdout, or other valid file path",
	)
	rootCmd.PersistentFlags().String(
		"log-level", "info",
		"set log level, support [debug, info, warn, error, panic, fatal]",
	)
	rootCmd.PersistentFlags().Bool(
		"log-disable", false,
		"diable logger output",
	)
	rootCmd.PersistentFlags().Bool(
		"log-hide-caller", false,
		"hide the logger caller info in the log",
	)
	rootCmd.PersistentFlags().Bool(
		"log-hide-time", false,
		"hide the logger time info in the log",
	)
	rootCmd.PersistentFlags().Bool(
		"log-hide-level", false,
		"hide the logger level info in the log",
	)
}

func Run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

//var vpFlag, vpConf *viper.Viper
var vpFlag *viper.Viper = viper.New()
var vpConf *viper.Viper = viper.New()

func PersistentPreRun(cmd *cobra.Command, args []string) {
	configHandler(cmd, args)
	logHandler(cmd, args)
}

func configHandler(cmd *cobra.Command, args []string) {
	log := zap.S()
	defer log.Sync()

	vpFlag.BindPFlags(cmd.Flags())

	var file *os.File
	var err error

	config := vpFlag.GetString("config")

	flag := cmd.Flags().Lookup("config")
	if flag.Changed {
		file, err = os.Open(config)
		if err != nil {
			log.Fatalw("open given config file fail", "file", config, "err", err)
		}
	} else {
		fn := func(config string) *os.File {
			file, err := os.Open(config)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				log.Fatalw("open default config file fail", "file", config, "err", err)
			}
			return file
		}
		cands := []string{config, "/etc/" + config}
		for _, cand := range cands {
			file = fn(cand)
			if file != nil {
				config = cand
				break
			}
		}
		if file == nil {
			vpFlag.Set("config", "")
			return
		}
	}
	defer file.Close()

	vpConf.SetConfigType("toml")
	err = vpConf.ReadConfig(file)
	if err != nil {
		log.Fatalw("may be config file is not invalid TOML format", "file", config, "err", err)
	}
}

func logHandler(cmd *cobra.Command, args []string) {
	if vpFlag.GetBool("log-disable") || vpConf.GetBool("log.disable") {
		zap.ReplaceGlobals(zap.NewNop())
		return
	}

	level := vpFlag.GetString("log-level")
	if !cmd.Flags().Lookup("log-level").Changed && vpConf.IsSet("log.level") {
		level = vpConf.GetString("log.level")
	}

	outpath := vpFlag.GetString("log-file")
	if !cmd.Flags().Lookup("log-file").Changed && vpConf.IsSet("log.file") {
		outpath = vpConf.GetString("log.file")
	}

	var hideLevel, hideCaller, hideTime bool

	if cmd.Flags().Lookup("log-hide-level").Changed {
		hideLevel = vpFlag.GetBool("log-hide-level")
	} else if vpConf.IsSet("log.hide-level") {
		hideLevel = vpConf.GetBool("log.hide-level")
	}

	if cmd.Flags().Lookup("log-hide-caller").Changed {
		hideCaller = vpFlag.GetBool("log-hide-caller")
	} else if vpConf.IsSet("log.hide-caller") {
		hideCaller = vpConf.GetBool("log.hide-caller")
	} else if level != "debug" {
		hideCaller = true
	}

	if cmd.Flags().Lookup("log-hide-time").Changed {
		hideTime = vpFlag.GetBool("log-hide-time")
	} else if vpConf.IsSet("log.hide-time") {
		hideTime = vpConf.GetBool("log.hide-time")
	}

	logger := zapLogger(level, outpath, hideLevel, hideCaller, hideTime)
	zap.ReplaceGlobals(logger)
}

func zapLogger(level, outpath string, hideLevel, hideCaller, hideTime bool) *zap.Logger {
	lvl := zap.InfoLevel
	switch level {
	case "debug":
		lvl = zap.DebugLevel
	case "info":
		lvl = zap.InfoLevel
	case "warn":
		lvl = zap.WarnLevel
	case "error":
		lvl = zap.ErrorLevel
	case "panic":
		lvl = zap.PanicLevel
	case "fatal":
		lvl = zap.FatalLevel
	}
	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		OutputPaths:      []string{outpath},
		ErrorOutputPaths: []string{outpath},

		//Encoding: "console",
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			MessageKey:    "M",
			StacktraceKey: "S",
			LineEnding:    "\n",

			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,

			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
			},
		},

		//Development:       true,
		//DisableCaller:     true,
		//DisableStacktrace: true,
		//Sampling: &SamplingConfig{
		//	Initial:    100,
		//	Thereafter: 100,
		//},
		//InitialFields: map[string]interface{}{
		//	"foo": "bar",
		//},
	}
	if outpath == "stdout" || outpath == "stderr" {
		zapCfg.DisableStacktrace = true
	}
	if hideCaller {
		zapCfg.DisableCaller = true
	}
	if hideTime {
		zapCfg.EncoderConfig.TimeKey = ""
	}
	if hideLevel {
		zapCfg.EncoderConfig.LevelKey = ""
	}

	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func initDefaultLogger() *zap.SugaredLogger {
	logger := zapLogger("info", "stderr", false, true, false)
	zap.ReplaceGlobals(logger)
	return zap.S()
}
