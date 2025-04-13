package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Run a restore task",
	GroupID: "command",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	defaultFs := LoadDefaultFlags(restoreCmd.Name())
	restoreFs := LoadRestoreFlags(restoreCmd.Name())

	restoreCmd.Flags().AddFlagSet(defaultFs)
	restoreCmd.Flags().AddFlagSet(restoreFs)

	restoreGroup := &cobra.Group{
		ID:    "service",
		Title: "Restore services:",
	}
	restoreCmd.AddGroup(restoreGroup)
}
