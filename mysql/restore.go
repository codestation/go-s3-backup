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

	"github.com/urfave/cli"
	"megpoid.xyz/go/go-s3-backup/common"
)

func Restore(c *cli.Context, filepath string) error {
	args := []string{
		"-h", c.String("host"),
		"-P", c.String("port"),
		"-u", c.String("user"),
	}

	if database := c.String("database"); database != "" {
		args = append(args, "-D", database)
	}

	if pass := c.String("password"); pass != "" {
		args = append(args, "-p", pass)
	}

	cmd := exec.Command(DumpApp, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()

	err := common.DecompressFileToInput(cmd, filepath)

	if err != nil {
		return fmt.Errorf("couldn't decompress file to input, %v", err)
	}

	return nil
}
