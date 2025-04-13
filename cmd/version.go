package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"go.megpoid.dev/go-s3-backup/version"
)

func printVersion() {
	slog.Info("GoApp",
		slog.String("version", version.Tag),
		slog.String("commit", version.Revision),
		slog.Time("date", version.LastCommit),
		slog.Bool("clean_build", !version.Modified),
	)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Prints the version and build info`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
