/*
Copyright 2018 codestation

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

package main

import (
	"github.com/urfave/cli"
)

var defaultFlags = []cli.Flag{
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

var backupFlags = []cli.Flag{
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

var restoreFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "schedule",
		Usage:  "cron schedule",
		Value:  "none",
		EnvVar: "SCHEDULE",
	},
	cli.StringFlag{
		Name:   "restore-file",
		Usage:  "restore from this file instead of searching for the most recent",
		EnvVar: "RESTORE_FILE",
	},
}

func backupCmd() cli.Command {
	name := "backup"
	return cli.Command{
		Name:  name,
		Usage: "run a backup task",
		Flags: append(defaultFlags, backupFlags...),
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
		Flags: append(defaultFlags, restoreFlags...),
		Subcommands: []cli.Command{
			gogsCmd(name),
			postgresCmd(name),
			mysqlCmd(name),
			tarballCmd(name),
		},
	}
}
