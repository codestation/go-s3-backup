package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.megpoid.dev/go-s3-backup/commands"
)

var restoreMysqlS3Cmd = &cobra.Command{
	Use:     "s3",
	Short:   "Connect to S3 store",
	GroupID: "store",
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("Run", "method", cmd.Parent().Parent().Name(), "service", cmd.Parent().Name(), "store", cmd.Name())
		return commands.RunTask(cmd.Parent().Parent().Name(), cmd.Parent().Name(), cmd.Name())
	},
}

func init() {
	restoreMysqlCmd.AddCommand(restoreMysqlS3Cmd)
	s3Fs := LoadS3Flags(restoreMysqlS3Cmd.Name())
	restoreMysqlS3Cmd.Flags().AddFlagSet(s3Fs)
}
