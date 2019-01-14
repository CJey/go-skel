package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"{=APPNAME=}/build"
)

func preRun(c *cobra.Command, args []string) {
}

func run(c *cobra.Command, args []string) {
	switch {
	case vpFlag.GetBool("build-indicator"):
		fmt.Println(build.BuildIndicator())
	case vpFlag.GetBool("build-hash"):
		fmt.Println(build.BuildHash())
	default:
		fmt.Print(build.Info())
	}
}
