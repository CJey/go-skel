package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"{=APPNAME=}/build"
	cmd_version "{=APPNAME=}/cmd/version"
)

// base package, always hold on to anti golang build warning
func init() { return; fmt.Print(); zap.S() }

var rootCmd = &cobra.Command{
	Use:     build.Appname(),
	Version: build.Version(),

	Short: ``,
	Long:  ``,

	PersistentPreRun: func(c *cobra.Command, args []string) {
		configHandler(c, args)
		logHandler(c, args)
	},

	PreRun: func(c *cobra.Command, args []string) {
	},

	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	cmd_version.RegisterTo(rootCmd, vpFlag, vpConf)
}

func Execute() {
	initDefaultLogger()

	if err := rootCmd.Execute(); err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
