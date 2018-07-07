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
	"log"
	"os"

	"github.com/urfave/cli"
	"megpoid.xyz/go/postgres-s3-backup/s3"
)

var build = "0" // build number set at compile-time
var saveDir = "/tmp"
var appPath = "/app/gogs/gogs"

func setBackupFunc(_ *cli.Context) error {
	s3.DoBackup = gogsBackup
	return nil
}

func setRestoreFunc(_ *cli.Context) error {
	s3.DoRestore = gogsRestore
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "gogs-s3-backup"
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Commands = []cli.Command{
		{
			Name:   "backup",
			Usage:  "run a backup task",
			Before: setBackupFunc,
			Action: s3.RunBackup,
			Flags:  append(s3.Flags, s3.BackupFlags...),
		},
		{
			Name:   "restore",
			Usage:  "run a restore task",
			Before: setRestoreFunc,
			Action: s3.RunRestore,
			Flags:  s3.Flags,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
