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

package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli"
	"megpoid.xyz/go/go-s3-backup/common"
)

func Backup(c *cli.Context) (string, error) {
	filepath := fmt.Sprintf("%s/mysql-backup-%s.sql.gz", SaveDir, time.Now().Format("20060102150405"))
	args := []string{
		"-h", c.String("host"),
		"-P", c.String("port"),
		"-u", c.String("user"),
	}

	if database := c.String("database"); database != "" {
		args = append(args, "-B", database)
	} else {
		args = append(args, "--all-databases")
	}

	if pass := c.String("password"); pass != "" {
		args = append(args, "-p", pass)
	}

	cmd := exec.Command(DumpApp, args...)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if c.Bool("compress") {
		err := common.CompressAppOutput(cmd, filepath)

		if err != nil {
			return "", fmt.Errorf("couldn't compress app output, %v", err)
		}
	} else {
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("couldn't execute %s, %v", DumpApp, err)
		}
	}

	return filepath, nil
}
