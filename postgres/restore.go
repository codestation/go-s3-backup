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

	"github.com/urfave/cli"
)

func Restore(c *cli.Context, filepath string) error {
	args := []string{
		"-h", c.String("host"),
		"-p", c.String("port"),
		"-U", c.String("user"),
	}

	if database := c.String("database"); database != "" {
		args = append(args, "-d", database)
	}

	args = append(args, filepath)

	env := os.Environ()

	if pass := c.String("password"); pass != "" {
		env = append(env, "PGPASSWORD="+pass)
	}

	cmd := exec.Command(RestoreApp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't execute pg_restore, %v", err)
	}

	return nil
}
