package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restorePostgresCmd = &cobra.Command{
	Use:     "postgres",
	Short:   "Connect to postgres service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	restoreCmd.AddCommand(restorePostgresCmd)

	databaseFs := LoadDatabaseFlags(restorePostgresCmd.Name())
	postgresFs := LoadPostgresFlags(restorePostgresCmd.Name())

	restorePostgresCmd.Flags().AddFlagSet(databaseFs)
	restorePostgresCmd.Flags().AddFlagSet(postgresFs)

	restorePostgresGroup := &cobra.Group{
		ID:    "store",
		Title: "Restore destinations:",
	}
	restorePostgresCmd.AddGroup(restorePostgresGroup)
}
