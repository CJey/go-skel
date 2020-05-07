package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/cjey/gbase"

	"go-skel/app"
)

// define root command
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
}

// execute root command
func Execute() {
	if err := _CMDRoot.Execute(); err != nil {
		os.Exit(1)
	}
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

	var _, err = gbase.ReplaceZapLogger(
		viper.GetString("log.level"),
		viper.GetString("log.file"),
		viper.GetString("log.encoding"),
		viper.GetBool("log.show-caller"),
	)
	if err != nil {
		zap.S().Fatalw("Replace zap logger fail", "err", err,
			"level", viper.GetString("log.level"),
			"file", viper.GetString("log.file"),
			"encoding", viper.GetString("log.encoding"),
			"show-caller", viper.GetBool("log.show-caller"),
		)
	}
}
