package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the version of the application
	Version = "v0.1.0"
	// GitCommit is the git commit hash
	GitCommit = "dev"
	// BuildDate is the date when the binary was built
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of l8k along with build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("l8k %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
