package version

import (
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the current version of this application",
	Long:  `Nothing fancy. Print the version of this application`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Polygon CLI Version %s\n", Version)
	},
}
