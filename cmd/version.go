/*
Copyright 2025 codestation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
