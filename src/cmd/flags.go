package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"{=APPNAME=}/build"
)

var vpFlag *viper.Viper = viper.New()
var vpConf *viper.Viper = viper.New()

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

func configHandler(c *cobra.Command, args []string) {
	log := zap.S()
	defer log.Sync()

	vpFlag.BindPFlags(c.Flags())

	var file *os.File
	var err error

	config := vpFlag.GetString("config")

	flag := c.Flags().Lookup("config")
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
