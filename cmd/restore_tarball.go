package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreTarballCmd = &cobra.Command{
	Use:     "tarball",
	Short:   "Connect to tarball service",
	GroupID: "service",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	restoreCmd.AddCommand(restoreTarballCmd)
	tarballFs := LoadTarballFlags(restoreTarballCmd.Name())
	restoreTarballCmd.Flags().AddFlagSet(tarballFs)

	restoreTarballGroup := &cobra.Group{
		ID:    "store",
		Title: "Restore destinations:",
	}
	restoreTarballCmd.AddGroup(restoreTarballGroup)
}
