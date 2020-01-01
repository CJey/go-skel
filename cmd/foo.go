package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-skel/context"
)

var _ = fmt.Print

var cmdFoo = &cobra.Command{
	Use:   `foo`,
	Run:   runFoo,
	Short: `Example cases`,
	Long:  ``,
}

func init() {
	cmd := cmdFoo
	supportConfigAndLogger(cmd)

	cmd.PersistentFlags().String("hello", "world!", "for testing")
	viper.BindPFlag("hello", cmd.PersistentFlags().Lookup("hello"))

	//cmdRoot.AddCommand(cmd)
}

func runFoo(cmd *cobra.Command, args []string) {
	handleConfigAndLogger(cmd)
	fmt.Printf("Hello %s\n\n", viper.GetString("hello"))

	ctx := context.Background("foo")
	ctx.L.Debugw("bar", "my level", "debug")
	ctx.L.Infow("bar", "my level", "info")
	ctx.L.Warnw("bar", "my level", "warn")
	ctx.L.Errorw("bar", "my level", "error")
	ctx1 := ctx.New("bar")
	ctx1.L.Debugw("bar", "my level", "debug")
	ctx1.L.Infow("bar", "my level", "info")
	ctx1.L.Warnw("bar", "my level", "warn")
	ctx1.L.Errorw("bar", "my level", "error")
}
