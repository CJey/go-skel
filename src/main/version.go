package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"{=APPNAME=}/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show full version infomation.",
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
	Run: RunVersion,
	//TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool("build-indicator", false, "show indicator that injected while building")
	versionCmd.Flags().MarkHidden("build-indicator")

	versionCmd.Flags().Bool("build-hash", false, "show build hash code")
	versionCmd.Flags().MarkHidden("build-hash")
}

func RunVersion(cmd *cobra.Command, args []string) {
	switch {
	case vpFlag.GetBool("build-indicator"):
		fmt.Println(build.BuildIndicator())
	case vpFlag.GetBool("build-hash"):
		fmt.Println(build.BuildHash())
	default:
		fmt.Print(build.Info())
	}
}
