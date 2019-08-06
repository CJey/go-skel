package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"{=APPNAME=}/build"
)

// base package, always hold on to anti golang build warning
func _() { fmt.Print(); zap.S() }

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
