package main

import (
	"github.com/urfave/cli"
)

var Flags = []cli.Flag{
	cli.IntFlag{
		Name:   "random-delay",
		Usage:  "schedule random delay",
		Value:  0,
		EnvVar: "SCHEDULE_RANDOM_DELAY",
	},
	cli.StringFlag{
		Name:   "savedir",
		Usage:  "directory to save/read backups",
		Value:  "/tmp",
		EnvVar: "SAVE_DIR",
	},
}

var BackupFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "schedule",
		Usage:  "cron schedule",
		Value:  "@daily",
		EnvVar: "SCHEDULE",
	},
	cli.IntFlag{
		Name:   "max-backups",
		Usage:  "max backups to keep (0 to disable the feature)",
		Value:  5,
		EnvVar: "MAX_BACKUPS",
	},
}

var RestoreFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "schedule",
		Usage:  "cron schedule",
		Value:  "none",
		EnvVar: "SCHEDULE",
	},
	cli.StringFlag{
		Name:   "s3key",
		Usage:  "s3 key",
		EnvVar: "S3_KEY",
	},
}

func backupCmd() cli.Command {
	name := "backup"
	return cli.Command{
		Name:  name,
		Usage: "run a backup task",
		Flags: append(Flags, BackupFlags...),
		Subcommands: []cli.Command{
			gogsCmd(name),
			postgresCmd(name),
			mysqlCmd(name),
			tarballCmd(name),
		},
	}
}

func restoreCmd() cli.Command {
	name := "restore"
	return cli.Command{
		Name:  "restore",
		Usage: "run a restore task",
		Flags: append(Flags, RestoreFlags...),
		Subcommands: []cli.Command{
			gogsCmd(name),
			postgresCmd(name),
			mysqlCmd(name),
			tarballCmd(name),
		},
	}
}
