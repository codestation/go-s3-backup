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
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/term"
	"megpoid.dev/go/go-s3-backup/version"
)

func printVersion(c *cli.Context) {
	slog.Info("go-s3-backup started",
		slog.String("version", version.Tag),
		slog.String("commit", version.Revision),
		slog.Time("date", version.LastCommit),
		slog.Bool("clean_build", !version.Modified),
	)
}

func main() {
	app := cli.NewApp()
	app.Usage = "run backups from various services to S3-like storage"
	app.Version = version.Tag
	cli.VersionPrinter = printVersion
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "load config from yaml file",
			EnvVars: []string{"CONFIG_FILE"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "enable debug logging",
			EnvVars: []string{"DEBUG"},
		},
	}

	app.Commands = []*cli.Command{
		backupCmd(),
		restoreCmd(),
	}

	app.Before = func(c *cli.Context) error {
		isTerminal := term.IsTerminal(int(os.Stdout.Fd()))
		if !isTerminal {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		}

		if c.Bool("debug") {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}

		slog.Info("go-s3-backup started",
			slog.String("version", version.Tag),
			slog.String("commit", version.Revision),
			slog.Time("date", version.LastCommit),
			slog.Bool("clean_build", !version.Modified),
		)

		if c.String("config") != "" {
			cfg, err := altsrc.NewYamlSourceFromFile(c.String("config"))
			if err != nil {
				app.Metadata = map[string]interface{}{
					"config": cfg,
				}
			}

			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Unrecoverable error", "error", err)
		os.Exit(1)
	}

	slog.Info("Shutdown complete")
}
