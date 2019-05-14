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
	"fmt"
	"os"

	log "gopkg.in/clog.v1"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

const versionFormatter = `go-s3-backup
Version:      %s
Git commit:   %s
Built:        %s
Compilation:  %s
`

func printVersion(c *cli.Context) {
	_, _ = fmt.Fprintf(c.App.Writer, versionFormatter, Version, Commit, BuildTime, BuildNumber)
}

func main() {
	app := cli.NewApp()
	app.Usage = "run backups from various services to S3-like storage"
	app.Version = Version
	cli.VersionPrinter = printVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "load config from yaml file",
			EnvVar: "CONFIG_FILE",
		},
	}

	app.Commands = []cli.Command{
		backupCmd(),
		restoreCmd(),
	}

	app.Before = func(c *cli.Context) error {
		if err := log.New(log.CONSOLE, log.ConsoleConfig{}); err != nil {
			return err
		}

		log.Info("Starting go-s3-backup, version: %s, commit: %s, built: %s, compilation: %s",
			Version,
			Commit,
			BuildTime,
			BuildNumber)

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
		log.Fatal(0, "Unrecoverable error: %v", err)
	}

	log.Info("Shutdown complete")
	log.Shutdown()
}
