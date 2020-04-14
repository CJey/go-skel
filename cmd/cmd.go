package cmd

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cjey/gbase"

	"go-skel/app"
)

var _CMDRoot = &cobra.Command{
	Use:   `{{appname}}`,
	Short: ``,
	Long:  ``,
}

func init() {
	var rpl = strings.NewReplacer(`{{appname}}`, app.Name)
	_CMDRoot.Use = rpl.Replace(_CMDRoot.Use)
	_CMDRoot.Short = rpl.Replace(_CMDRoot.Short)
	_CMDRoot.Long = rpl.Replace(_CMDRoot.Long)
	_CMDRoot.Version = app.Version
	supportConfigAndLogger(_CMDRoot)

	cobra.OnInitialize(commandInitializer)
}

func Execute() {
	if err := _CMDRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

func commandInitializer() {
	// initialize random generator
	rand.Seed(gbase.BootTime.UnixNano())

	// initialize default logger
	setLogger("debug", "stderr", "console", false)
}

func supportConfigAndLogger(cmd *cobra.Command) {
	supportConfig(cmd)
	supportLogger(cmd)
}

// used when command really run
func handleConfigAndLogger(cmd *cobra.Command) {
	handleConfig(cmd)
	handleLogger(cmd)
}

func supportConfig(cmd *cobra.Command) {
	// do not bind, flag only
	// -c, --config ${app.Name}.toml
	cmd.PersistentFlags().StringP("config", "c", app.Name+".toml",
		"configuration file path")
}

func handleConfig(cmd *cobra.Command) {
	var fpath, _ = cmd.Flags().GetString("config")

	viper.SetConfigFile(fpath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		if cmd.Flag("config").Changed == false {
			// use default config file
			// so, not found error must be ignored
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return
			}
			if errors.Is(err, os.ErrNotExist) {
				return
			}
		}
		zap.S().Fatalw("Read configuration file fail", "err", err, "file", fpath)
	}
}

func supportLogger(cmd *cobra.Command) {
	// --log-level debug
	cmd.PersistentFlags().String("log-level", "debug",
		"set log level, support [debug, info, warn, error, panic, fatal]")
	// --log-file stderr
	cmd.PersistentFlags().String("log-file", "stderr",
		"log file path, support stderr, stdout, or other valid file path")
	// --log-encoding console
	cmd.PersistentFlags().String("log-encoding", "console",
		"set log encoding, support [console, json]")
	// --log-show-caller=false
	cmd.PersistentFlags().Bool("log-show-caller", false,
		"show the logger caller info in the log")
}

func handleLogger(cmd *cobra.Command) {
	// bind flags
	viper.BindPFlag("log.level", cmd.Flags().Lookup("log-level"))
	viper.BindPFlag("log.file", cmd.Flags().Lookup("log-file"))
	viper.BindPFlag("log.encoding", cmd.Flags().Lookup("log-encoding"))
	viper.BindPFlag("log.show-caller", cmd.Flags().Lookup("log-show-caller"))

	setLogger(
		viper.GetString("log.level"),
		viper.GetString("log.file"),
		viper.GetString("log.encoding"),
		viper.GetBool("log.show-caller"),
	)
}

func setLogger(lvl, outpath, encoding string, showCaller bool) {
	var level zapcore.Level
	switch lvl {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		zap.S().Fatalw("unexpected log level", "level", lvl)
	}
	l, err := logger(level, outpath, encoding, showCaller)
	if err != nil {
		zap.S().Fatalw("build logger fail", "err", err, "file", outpath)
	}
	zap.ReplaceGlobals(l)
}

func logger(level zapcore.Level, outpath, encoding string, showCaller bool) (*zap.Logger, error) {
	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{outpath},
		ErrorOutputPaths: []string{outpath},

		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "T",
			LevelKey:   "L",
			NameKey:    "N",
			MessageKey: "M",
			LineEnding: "\n",

			CallerKey:     "Caller",
			StacktraceKey: "Stack",

			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
			},
		},

		//Development:       true,
		DisableCaller: !showCaller,
		//DisableStacktrace: true,
		//Sampling: &SamplingConfig{
		//  Initial:    100,
		//  Thereafter: 100,
		//},
		//InitialFields: map[string]interface{}{
		//  "foo": "bar",
		//},
	}

	return zapCfg.Build()
}
