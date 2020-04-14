package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cjey/gbase"
)

var _CMDFoo = &cobra.Command{
	Use:   `foo`,
	Run:   runFoo,
	Short: `Example cases`,
	Long:  ``,
}

func init() {
	var cmd = _CMDFoo
	supportConfigAndLogger(cmd)

	// --hello world!
	cmd.PersistentFlags().String("hello", "world!", "for testing")

	_CMDRoot.AddCommand(cmd)
}

func runFoo(cmd *cobra.Command, args []string) {
	handleConfigAndLogger(cmd)

	// bind flags
	viper.BindPFlag("hello", cmd.PersistentFlags().Lookup("hello"))

	fmt.Printf("Hello %s\n\n", viper.GetString("hello"))

	var ctx = gbase.NamedContext("foo")
	ctx.Debug("bar", "my level", "debug")
	ctx.Info("bar", "my level", "info")
	ctx.Warn("bar", "my level", "warn")
	ctx.Error("bar", "my level", "error")

	var ctx1 = ctx.ForkAt("bar")
	ctx1.Debug("bar", "my level", "debug")
	ctx1.Info("bar", "my level", "info")
	ctx1.Warn("bar", "my level", "warn")
	ctx1.Error("bar", "my level", "error")
}
