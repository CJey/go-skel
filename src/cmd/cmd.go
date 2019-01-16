package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"{=APPNAME=}/build"
	cmd_version "{=APPNAME=}/cmd/version"
	cmd_hello "{=APPNAME=}/cmd/hello"
)

var rootCmd = &cobra.Command{
	Use:     build.Appname(),
	Version: build.Version(),

	Short: ``,
	Long:  ``,

	PersistentPreRun: func(c *cobra.Command, args []string) {
        // global handlers
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
	cmd_hello.RegisterTo(rootCmd, vpFlag, vpConf)
}

func Execute() {
	initDefaultLogger()

	if err := rootCmd.Execute(); err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
