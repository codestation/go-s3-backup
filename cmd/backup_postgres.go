package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupPostgresCmd = &cobra.Command{
	Use:     "postgres",
	Short:   "Connect to postgres service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	backupCmd.AddCommand(backupPostgresCmd)

	databaseFs := LoadDatabaseFlags(backupPostgresCmd.Name())
	postgresFs := LoadPostgresFlags(backupPostgresCmd.Name())

	backupPostgresCmd.Flags().AddFlagSet(databaseFs)
	backupPostgresCmd.Flags().AddFlagSet(postgresFs)

	backupPostgresGroup := &cobra.Group{
		ID:    "store",
		Title: "Backup destinations:",
	}
	backupPostgresCmd.AddGroup(backupPostgresGroup)
}
