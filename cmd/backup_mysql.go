package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupMysqlCmd = &cobra.Command{
	Use:     "mysql",
	Short:   "Connect to mysql/mariadb service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	backupCmd.AddCommand(backupMysqlCmd)

	databaseFs := LoadDatabaseFlags(backupMysqlCmd.Name())
	mysqlFs := LoadMySQLFlags(backupMysqlCmd.Name())

	backupMysqlCmd.Flags().AddFlagSet(databaseFs)
	backupMysqlCmd.Flags().AddFlagSet(mysqlFs)

	backupMysqlGroup := &cobra.Group{
		ID:    "store",
		Title: "Backup destinations:",
	}
	backupMysqlCmd.AddGroup(backupMysqlGroup)
}
