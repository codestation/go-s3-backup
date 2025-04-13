package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreMysqlCmd = &cobra.Command{
	Use:     "mysql",
	Short:   "Connect to mysql/mariadb service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	restoreCmd.AddCommand(restoreMysqlCmd)

	databaseFs := LoadDatabaseFlags(restoreMysqlCmd.Name())
	mysqlFs := LoadMySQLFlags(restoreMysqlCmd.Name())

	restoreMysqlCmd.Flags().AddFlagSet(databaseFs)
	restoreMysqlCmd.Flags().AddFlagSet(mysqlFs)

	restoreMysqlGroup := &cobra.Group{
		ID:    "store",
		Title: "Restore destinations:",
	}
	restoreMysqlCmd.AddGroup(restoreMysqlGroup)
}
