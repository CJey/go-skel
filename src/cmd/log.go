package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initDefaultLogger() *zap.SugaredLogger {
	logger := zapLogger("info", "stderr", false, true, false)
	zap.ReplaceGlobals(logger)
	return zap.S()
}

func logHandler(c *cobra.Command, args []string) {
	if vpFlag.GetBool("log-disable") || vpConf.GetBool("log.disable") {
		zap.ReplaceGlobals(zap.NewNop())
		return
	}

	level := vpFlag.GetString("log-level")
	if !c.Flags().Lookup("log-level").Changed && vpConf.IsSet("log.level") {
		level = vpConf.GetString("log.level")
	}

	outpath := vpFlag.GetString("log-file")
	if !c.Flags().Lookup("log-file").Changed && vpConf.IsSet("log.file") {
		outpath = vpConf.GetString("log.file")
	}

	var hideLevel, hideCaller, hideTime bool

	if c.Flags().Lookup("log-hide-level").Changed {
		hideLevel = vpFlag.GetBool("log-hide-level")
	} else if vpConf.IsSet("log.hide-level") {
		hideLevel = vpConf.GetBool("log.hide-level")
	}

	if c.Flags().Lookup("log-hide-caller").Changed {
		hideCaller = vpFlag.GetBool("log-hide-caller")
	} else if vpConf.IsSet("log.hide-caller") {
		hideCaller = vpConf.GetBool("log.hide-caller")
	} else if level != "debug" {
		hideCaller = true
	}

	if c.Flags().Lookup("log-hide-time").Changed {
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
