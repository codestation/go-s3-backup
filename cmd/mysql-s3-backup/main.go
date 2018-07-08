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

	"megpoid.xyz/go/go-s3-backup/mysql"
	"megpoid.xyz/go/go-s3-backup/s3"

	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func setBackupFunc(_ *cli.Context) error {
	s3.DoBackup = mysql.Backup
	return nil
}

func setRestoreFunc(_ *cli.Context) error {
	s3.DoRestore = mysql.Restore
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "mysql-s3-backup"
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
			Flags:  append(s3.Flags, s3.RestoreFlags...),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
