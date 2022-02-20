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
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "unknwon.dev/clog/v2"
)

// MySQLConfig has the config options for the MySQLConfig service
type MySQLConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	NamePrefix     string
	NameAsPrefix   bool
	Options        string
	Compress       bool
	SaveDir        string
	IgnoreExitCode bool
}

// MysqlDumpApp points to the mysqldump binary location
var MysqlDumpApp = "/usr/bin/mysqldump"

// MysqlRestoreApp points to the mysql binary location
var MysqlRestoreApp = "/usr/bin/mysql"

func (m *MySQLConfig) newBaseArgs() []string {
	args := []string{
		"-h", m.Host,
		"-P", m.Port,
		"-u", m.User,
	}

	if m.Password != "" {
		args = append(args, "-p"+m.Password)
	}

	options := strings.Fields(m.Options)

	// add extra options
	if len(options) > 0 {
		args = append(args, options...)
	}

	return args
}

func (m *MySQLConfig) getNamePrefix() string {
	var prefix string
	if m.NameAsPrefix {
		prefix = m.Database
	} else if m.NamePrefix != "" {
		prefix = m.NamePrefix
	} else {
		prefix = "mysql-backup"
	}

	return prefix
}

// Backup generates a dump of the database and returns the path where is stored
func (m *MySQLConfig) Backup() (*BackupResults, error) {
	namePrefix := m.getNamePrefix()
	filepath := generateFilename(m.SaveDir, namePrefix)
	args := m.newBaseArgs()

	if m.Database != "" {
		args = append(args, "-B", m.Database)
	} else {
		args = append(args, "--all-databases")
	}

	if !m.Compress {
		filepath += ".sql"
		args = append(args, "-r", filepath)
	} else {
		filepath += ".sql.gz"
	}

	app := CmdConfig{CensorArg: "-p"}

	if err := os.MkdirAll(m.SaveDir, 0755); err != nil {
		return nil, err
	}

	if m.Compress {
		f, err := os.Create(filepath)
		if err != nil {
			return nil, fmt.Errorf("cannot create file: %v", err)
		}

		defer f.Close()

		writer := gzip.NewWriter(f)
		defer writer.Close()

		app.OutputFile = writer
	}

	if err := app.CmdRun(MysqlDumpApp, args...); err != nil {
		return nil, fmt.Errorf("couldn't execute %s, %v", MysqlDumpApp, err)
	}

	result := &BackupResults{[]BackupResult{
		{
			NamePrefix: namePrefix,
			DirPrefix:  m.Database,
			Path:       filepath,
		},
	}}

	return result, nil
}

// Restore takes a database dump and restores it
func (m *MySQLConfig) Restore(filepath string) error {
	args := m.newBaseArgs()
	app := CmdConfig{}

	if m.Database != "" {
		args = append(args, "-D", m.Database)
	}

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}

	defer f.Close()

	if strings.HasSuffix(filepath, ".gz") {
		reader, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("cannot create gzip reader: %v", err)
		}

		defer reader.Close()
		app.InputFile = reader
	} else {
		app.InputFile = f
	}

	if err := app.CmdRun(MysqlRestoreApp, args...); err != nil {
		serr, ok := err.(*exec.ExitError)

		if ok && m.IgnoreExitCode {
			log.Info("Ignored exit code of restore process: %v", serr)
		} else {
			return fmt.Errorf("couldn't execute %s, %v", MysqlRestoreApp, err)
		}
	}

	return nil
}
