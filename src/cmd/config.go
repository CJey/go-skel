package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var vpConf *viper.Viper = viper.New()

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
