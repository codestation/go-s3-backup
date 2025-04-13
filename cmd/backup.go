package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Run a backup task",
	GroupID: "command",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	defaultFs := LoadDefaultFlags(backupCmd.Name())
	backupFs := LoadBackupFlags(backupCmd.Name())

	backupCmd.Flags().AddFlagSet(defaultFs)
	backupCmd.Flags().AddFlagSet(backupFs)

	backupGroup := &cobra.Group{
		ID:    "service",
		Title: "Backup services:",
	}
	backupCmd.AddGroup(backupGroup)
}
