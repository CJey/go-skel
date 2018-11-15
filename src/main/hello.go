package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use: "hello",
	Run: RunHello,
	//TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(helloCmd)
}

func RunHello(cmd *cobra.Command, args []string) {
	fmt.Println("world!")
}
