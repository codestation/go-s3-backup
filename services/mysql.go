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
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	log "unknwon.dev/clog/v2"
)

// MySQLConfig has the config options for the MySQLConfig service
type MySQLConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	Database         string
	NamePrefix       string
	NameAsPrefix     bool
	Options          string
	Compress         bool
	SaveDir          string
	SplitDatabases   bool
	ExcludeDatabases []string
	IgnoreExitCode   bool
}

// MysqlDumpApp points to the mysqldump binary location
var MysqlDumpApp = "/usr/bin/mysqldump"

// MysqlCmdApp points to the mysql binary location
var MysqlCmdApp = "/usr/bin/mysql"

var mysqlListDatabasesQuery = "show databases"

func (m *MySQLConfig) newBaseArgs(skipOptions bool) []string {
	args := []string{
		"-h", m.Host,
		"-P", m.Port,
		"-u", m.User,
	}

	if m.Password != "" {
		args = append(args, "-p"+m.Password)
	}

	if !skipOptions {
		options := strings.Fields(m.Options)

		// add extra options
		if len(options) > 0 {
			args = append(args, options...)
		}
	}

	return args
}

func (m *MySQLConfig) getNamePrefix() string {
	var prefix string
	switch {
	case m.NameAsPrefix && m.Database != "":
		prefix = m.Database
	case m.NamePrefix != "":
		prefix = m.NamePrefix
	default:
		prefix = "mysql-backup"
	}

	return prefix
}

// Backup generates a dump of the database and returns the path where is stored
func (m *MySQLConfig) Backup() (*BackupResults, error) {
	switch {
	case m.SplitDatabases:
		databases, err := m.listDatabases()
		if err != nil {
			return nil, fmt.Errorf("failed to list databases, %w", err)
		}

		var resultList []BackupResult

		for _, database := range databases {

			if len(m.ExcludeDatabases) > 0 {
				found := false
				for _, exclude := range m.ExcludeDatabases {
					matched, err := path.Match(exclude, database)
					if err != nil {
						log.Error("Invalid pattern %s, skipping", exclude)
						found = true
						break
					}
					if matched {
						log.Info("Excluding database %s that match excluded pattern %s", database, exclude)
						found = true
						break
					}
				}
				if found {
					continue
				}
			}

			m.Database = database
			namePrefix := m.getNamePrefix()
			filepath, err := m.backupDatabase("", namePrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to backup database %s, %w", database, err)
			}
			resultEntry := BackupResult{DirPrefix: namePrefix, NamePrefix: namePrefix, Path: filepath}
			resultList = append(resultList, resultEntry)
		}

		result := &BackupResults{resultList}
		return result, nil
	default:
		namePrefix := m.getNamePrefix()
		filepath, err := m.backupDatabase("", namePrefix)
		if err != nil {
			return nil, err
		}

		return &BackupResults{Entries: []BackupResult{{
			NamePrefix: namePrefix,
			Path:       filepath,
		}}}, nil
	}
}

func (m *MySQLConfig) backupDatabase(basedir, namePrefix string) (string, error) {
	savePath := path.Join(m.SaveDir, basedir)
	filepath := generateFilename(savePath, namePrefix)
	args := m.newBaseArgs(false)

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

	if err := os.MkdirAll(m.SaveDir, 0o755); err != nil {
		return "", err
	}

	if m.Compress {
		f, err := os.Create(filepath)
		if err != nil {
			return "", fmt.Errorf("cannot create file: %v", err)
		}

		defer f.Close()

		writer := gzip.NewWriter(f)
		defer writer.Close()

		app.OutputFile = writer
	}

	if err := app.CmdRun(MysqlDumpApp, args...); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", MysqlDumpApp, err)
	}

	return filepath, nil
}

// Restore takes a database dump and restores it
func (m *MySQLConfig) Restore(filepath string) error {
	args := m.newBaseArgs(false)
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

	if err := app.CmdRun(MysqlCmdApp, args...); err != nil {
		serr, ok := err.(*exec.ExitError)

		if ok && m.IgnoreExitCode {
			log.Info("Ignored exit code of restore process: %v", serr)
		} else {
			return fmt.Errorf("couldn't execute %s, %v", MysqlCmdApp, err)
		}
	}

	return nil
}

func (m *MySQLConfig) listDatabases() ([]string, error) {
	args := m.newBaseArgs(true)
	args = append(args, "-s", "--skip-column-names", "-r")
	app := CmdConfig{}

	var b bytes.Buffer
	outputWriter := bufio.NewWriter(&b)
	app.OutputFile = outputWriter

	var listDatabases []string
	listDatabases = append(listDatabases, args...)
	listDatabases = append(listDatabases, "-e", mysqlListDatabasesQuery)

	if err := app.CmdRun(MysqlCmdApp, listDatabases...); err != nil {
		return nil, fmt.Errorf("mysqldump error on database list, %w", err)
	}

	scanner := bufio.NewScanner(&b)
	var databases []string
	for scanner.Scan() {
		databases = append(databases, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse psql output, %w", err)
	}

	return databases, nil
}
