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
	"os"

	"megpoid.xyz/go/go-s3-backup/cmd"
	"megpoid.xyz/go/go-s3-backup/services"

	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
)

var (
	// Build information (set by -ldflags)
	BuildTime   string
	BuildCommit string
	AppVersion  string
)

func getService(c *cli.Context) services.Service {
	var serv services.Service
	switch c.Command.Name {
	case "gogs":
		serv = cmd.NewGogsConfig(c)
	case "mysql":
		serv = cmd.NewMysqlConfig(c)
	case "postgres":
		serv = cmd.NewPostgresConfig(c)
	default:
		log.Fatal(0, "unsupported service: %s", c.Args().Get(1))
	}

	return serv
}

func backupJob(c *cli.Context) error {
	return cmd.RunScheduler(c, func(c *cli.Context) error {
		serv := getService(c)
		s3 := cmd.NewS3Config(c)
		return cmd.BackupTask(c, serv, s3)
	})
}

func restoreJob(c *cli.Context) error {
	return cmd.RunScheduler(c, func(c *cli.Context) error {
		serv := getService(c)
		s3 := cmd.NewS3Config(c)
		return cmd.RestoreTask(c, serv, s3)
	})
}

func main() {
	app := cli.NewApp()
	app.Usage = "run backups from various services to S3-like storage"
	app.Version = AppVersion
	app.Commands = []cli.Command{
		{
			Name:  "backup",
			Usage: "run a backup task",
			Flags: append(cmd.Flags, cmd.BackupFlags...),
			Subcommands: []cli.Command{
				{
					Name:   "gogs",
					Action: backupJob,
					Flags:  cmd.GogsFlags,
				},
				{
					Name:   "mysql",
					Action: backupJob,
					Flags:  cmd.DatabaseFlags,
				},
				{
					Name:   "postgres",
					Action: backupJob,
					Flags:  append(cmd.DatabaseFlags, cmd.PostgresFlags...),
				},
			},
		},
		{
			Name:  "restore",
			Usage: "run a restore task",
			Flags: append(cmd.Flags, cmd.RestoreFlags...),
			Subcommands: []cli.Command{
				{
					Name:   "gogs",
					Action: restoreJob,
					Flags:  cmd.GogsFlags,
				},
				{
					Name:   "mysql",
					Action: restoreJob,
					Flags:  cmd.DatabaseFlags,
				},
				{
					Name:   "postgres",
					Action: restoreJob,
					Flags:  append(cmd.DatabaseFlags, cmd.PostgresFlags...),
				},
			},
		},
	}

	log.New(log.CONSOLE, log.ConsoleConfig{})
	log.Info("go-s3-backup %s", AppVersion)

	if len(BuildTime) > 0 {
		log.Trace("Build Time: %s", BuildTime)
		log.Trace("Build Commit: %s", BuildCommit)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(0, err.Error())
	}
}
