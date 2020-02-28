package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-skel/app"
)

var cmdVersion = &cobra.Command{
	Use:   `version`,
	Run:   runVersion,
	Short: `Show full version infomation.`,
	Long: `Show full version infomation.
Include git infomation & build infomation.

Specially, "GitTrace" format:

    {gitNumber}.{gitShortHash}
    or (when build at unclean git working directory)
    {gitNumber}.{gitShortHash} # {gitStatusNumber}.{gitStatusHash}

gitNumber: how many commits since first commit until current commit
gitShortHash: 7 chars at current commit hash code's head
gitStatusNumber: how many different files (or untracked by git) compare to current commit
gitStatusHash: 7 chars at the hash code's head which indicate different
`,
}

func init() {
	cmd := cmdVersion

	cmd.PersistentFlags().Bool("tag", false, "Show info about git tag")

	viper.BindPFlag("tag", cmd.PersistentFlags().Lookup("tag"))

	cmdRoot.AddCommand(cmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	var tpl *template.Template
	if viper.GetBool("tag") {
		if app.App().Git.TagName == "" {
			return
		}
		tpl = template.Must(template.New("info").Parse(`
AppName     {{.Name}}
TagName     {{.Git.TagName}}
{{if .Git.TagTrace}}
TagTrace    {{.Git.TagTrace}}{{if .Git.Branch}}
GitBranch   {{.Git.Branch}}{{end}}{{if .Git.Repo}}
GitRepo     {{.Git.Repo}}{{end}}
TagHash     {{.Git.TagHash}} @ {{.Git.TagTimeString}}{{end}}
---

{{.Git.TagMessage}}
`))
	} else {
		tpl = template.Must(template.New("info").Parse(`
AppName     {{.Name}}
Version     {{.Version}}-{{.Release}}
{{if .Git.Trace}}
GitTrace    {{.Git.Trace}}{{if .Git.Branch}}
GitBranch   {{.Git.Branch}}{{end}}{{if .Git.Repo}}
GitRepo     {{.Git.Repo}}{{end}}
GitHash     {{.Git.CommitHash}} @ {{.Git.CommitTimeString}}
{{end}}
Golang      {{.Go.Version}} {{.Go.Arch}}
BuildInfo   {{.Build.ID}} @ {{.Build.TimeString}}
`))
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, app.App()); err != nil {
		panic(err)
	}
	buf.Next(1) // trim first \n

	fmt.Print(buf.String())
}
