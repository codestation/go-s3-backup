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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreMysqlCmd = &cobra.Command{
	Use:     "mysql",
	Short:   "Connect to mysql/mariadb service",
	GroupID: "service",
	Aliases: []string{"mariadb"},
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	restoreCmd.AddCommand(restoreMysqlCmd)

	databaseFs := LoadDatabaseFlags(restoreMysqlCmd.Name())
	mysqlFs := LoadMySQLFlags(restoreMysqlCmd.Name())

	restoreMysqlCmd.PersistentFlags().AddFlagSet(databaseFs)
	restoreMysqlCmd.PersistentFlags().AddFlagSet(mysqlFs)

	restoreMysqlGroup := &cobra.Group{
		ID:    "store",
		Title: "Restore destinations:",
	}
	restoreMysqlCmd.AddGroup(restoreMysqlGroup)
}
