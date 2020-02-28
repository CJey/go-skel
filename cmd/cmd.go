package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go-skel/app"
)

var _ = fmt.Print

func supportConfigAndLogger(cmd *cobra.Command) {
	supportConfig(cmd)
	supportLogger(cmd)
}

func handleConfigAndLogger(cmd *cobra.Command) {
	handleConfig(cmd)
	handleLogger(cmd)
}

func supportConfig(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "c", app.Name+".toml",
		"configuration file path")
}
func handleConfig(cmd *cobra.Command) {
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))

	fpath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}
	if fpath == "" {
		return
	}

	viper.SetConfigFile(fpath)
	viper.SetConfigType("toml")
	if err = viper.ReadInConfig(); err != nil {
		if !cmd.Flags().Lookup("config").Changed {
			// default config file not found, ignored
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return
			}
			if errors.Is(err, os.ErrNotExist) {
				return
			}
		}
		zap.S().Fatalw("read configuration file fail", "err", err, "file", fpath)
	}
	return
}

func supportLogger(cmd *cobra.Command) {
	cmd.PersistentFlags().String("log-level", "debug",
		"set log level, support [debug, info, warn, error, panic, fatal]")
	cmd.PersistentFlags().String("log-file", "stderr",
		"log file path, support stderr, stdout, or other valid file path")
	cmd.PersistentFlags().String("log-encoding", "console",
		"set log encoding, support [console, json]")
	cmd.PersistentFlags().Bool("log-show-caller", false,
		"show the logger caller info in the log")
}

func handleLogger(cmd *cobra.Command) {
	viper.BindPFlag("log.level", cmd.Flags().Lookup("log-level"))
	viper.BindPFlag("log.file", cmd.Flags().Lookup("log-file"))
	viper.BindPFlag("log.encoding", cmd.Flags().Lookup("log-encoding"))
	viper.BindPFlag("log.show-caller", cmd.Flags().Lookup("log-show-caller"))

	lvl := viper.GetString("log.level")
	f := viper.GetString("log.file")
	enc := viper.GetString("log.encoding")
	show := viper.GetBool("log.show-caller")

	setLogger(lvl, f, enc, show)
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
