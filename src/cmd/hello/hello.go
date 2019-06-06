package hello

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vpFlag *viper.Viper
var vpConf *viper.Viper

var cmd = &cobra.Command{
	Use:    "hello",
	Short:  ``,
	Long:   ``,
	PreRun: preRun,
	Run:    run,
}

func init() {
}

func RegisterTo(father *cobra.Command, flag, conf *viper.Viper) {
	father.AddCommand(cmd)
	vpFlag = flag
	vpConf = conf
}

func Called(c *cobra.Command) bool {
	return c == cmd
}
