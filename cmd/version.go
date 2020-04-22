package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"go-skel/app"
)

var _CMDVersion = &cobra.Command{
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
	var cmd = _CMDVersion

	// do not bind, flag only
	// --changelog=false
	cmd.PersistentFlags().Bool("changelog", false, "Show info about git tag history")
	// --max-lines 32
	cmd.PersistentFlags().Uint("max-lines", 32, "If show changelog, the max lines to show, 0 means unlimit")

	_CMDRoot.AddCommand(cmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	if showChangelog, _ := cmd.Flags().GetBool("changelog"); showChangelog {
		var (
			changelog  = strings.TrimSpace(app.App().Git.TagMessage)
			maxline, _ = cmd.Flags().GetUint("max-lines")
			lines      = strings.Split(changelog, "\n")
		)
		if maxline > 0 && maxline < uint(len(lines)) {
			lines = lines[:maxline]
			changelog = strings.Join(lines, "\n")
		}
		fmt.Println(changelog)
		return
	}

	var tpl = template.Must(template.New("info").Parse(`
AppName     {{.Name}}
Version     {{.FullVersion}}
{{if .Git.Trace}}
GitTrace    {{.Git.Trace}}{{if .Git.Branch}}
GitBranch   {{.Git.Branch}}{{end}}{{if .Git.Repo}}
GitRepo     {{.Git.Repo}}{{end}}
GitHash     {{.Git.CommitHash}} @ {{.Git.CommitTimeString}}
{{end}}
Golang      {{.Go.Version}} {{.Go.Arch}}
BuildInfo   {{.Build.ID}} @ {{.Build.TimeString}}
`))

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, app.App()); err != nil {
		panic(err)
	}
	buf.Next(1) // trim first \n

	fmt.Print(buf.String())
}
