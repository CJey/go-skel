package cmd

import (
	"math/rand"
	"strings"

	"github.com/spf13/cobra"

	"go-skel/app"
	"go-skel/context"
)

var cmdRoot = &cobra.Command{
	Use:   `{{appname}}`,
	Short: ``,
	Long:  ``,
}

func init() {
	rpl := strings.NewReplacer(`{{appname}}`, app.Name)
	cmdRoot.Use = rpl.Replace(cmdRoot.Use)
	cmdRoot.Short = rpl.Replace(cmdRoot.Short)
	cmdRoot.Long = rpl.Replace(cmdRoot.Long)
	cmdRoot.Version = app.Version
	supportConfigAndLogger(cmdRoot)

	cobra.OnInitialize(commandInitializer)
}

func commandInitializer() {
	// initialize random generator
	rand.Seed(app.BootTime.UnixNano())

	// initialize default logger
	setLogger("debug", "stderr", "console", false)

	// initialize context
	context.BootID(app.BootID)
}
