package cmd

import (
	"os"
)

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
