package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupTarballCmd = &cobra.Command{
	Use:     "tarball",
	Short:   "Connect to tarball service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	backupCmd.AddCommand(backupTarballCmd)
	tarballFs := LoadMySQLFlags(backupTarballCmd.Name())
	backupTarballCmd.Flags().AddFlagSet(tarballFs)

	backupTarballGroup := &cobra.Group{
		ID:    "store",
		Title: "Backup destinations:",
	}
	backupTarballCmd.AddGroup(backupTarballGroup)
}
