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

package services

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type MySQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  []string
	Compress bool
}

var MysqlDumpApp = "/usr/bin/mysqldump"
var MysqlRestoreApp = "/usr/bin/mysql"

func (m MySQL) Backup() (string, error) {
	filepath := fmt.Sprintf("%s/mysql-backup-%s", SaveDir, time.Now().Format("20060102150405"))

	args := []string{
		"-h", m.Host,
		"-P", m.Port,
		"-u", m.User,
	}

	if m.Database != "" {
		args = append(args, "-B", m.Database)
	} else {
		args = append(args, "--all-databases")
	}

	if m.Password != "" {
		args = append(args, "-p", m.Password)
	}

	if !m.Compress {
		filepath += ".sql"
		args = append(args, "-r", filepath)
	} else {
		filepath += ".sql.gz"
	}

	// add extra options
	if len(m.Options) > 0 {
		args = append(args, m.Options...)
	}

	cmd := exec.Command(MysqlDumpApp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if m.Compress {
		if err := CompressAppOutput(cmd, filepath); err != nil {
			return "", fmt.Errorf("couldn't compress app output, %v", err)
		}
	} else {
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("couldn't run mysql dump, %v", err)
		}
	}

	return filepath, nil
}

func (m *MySQL) Restore(filepath string) error {
	args := []string{
		"-h", m.Host,
		"-P", m.Port,
		"-u", m.User,
	}

	if m.Database != "" {
		args = append(args, "-B", m.Database)
	} else {
		args = append(args, "--all-databases")
	}

	if m.Password != "" {
		args = append(args, "-p", m.Password)
	}

	// add extra options
	if len(m.Options) > 0 {
		args = append(args, m.Options...)
	}

	cmd := exec.Command(MysqlRestoreApp, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()

	if err := ReadFileToInput(cmd, filepath); err != nil {
		return fmt.Errorf("couldn't decompress file to input, %v", err)
	}

	return nil
}
