package cmd

import (
	"github.com/spf13/viper"

	"{=APPNAME=}/build"
)

var vpFlag *viper.Viper = viper.New()

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
