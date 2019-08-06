package hello

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// base package, always hold on to anti golang build warning
func _() { fmt.Print(); zap.S() }

func preRun(c *cobra.Command, args []string) {
}

func run(c *cobra.Command, args []string) {
	fmt.Println("word!")
}
