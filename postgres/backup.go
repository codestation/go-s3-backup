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

package postgres

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli"
)

func Backup(c *cli.Context) (string, error) {
	filepath := fmt.Sprintf("%s/postgres-backup-%s.dump", SaveDir, time.Now().Format("20060102150405"))
	args := []string{
		"-h", c.String("host"),
		"-p", c.String("port"),
		"-U", c.String("user"),
		"-f", filepath,
	}

	var app string

	if database := c.String("database"); database != "" {
		args = append(args, "-d", database)
		app = DumpApp
	} else {
		app = DumpallApp
	}

	env := os.Environ()

	if pass := c.String("password"); pass != "" {
		env = append(env, "PGPASSWORD="+pass)
	}

	cmd := exec.Command(app, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("couldn't execute pg_dump, %v", err)
	}

	return filepath, nil
}
