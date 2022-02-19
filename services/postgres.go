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

// PostgresConfig has the config options for the PostgresConfig service
type PostgresConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	NamePrefix     string
	NameAsPrefix   bool
	Options        string
	Compress       bool
	Custom         bool
	SaveDir        string
	IgnoreExitCode bool
	Drop           bool
	Owner          string
}

// PostgresDumpApp points to the pg_dump binary location
var PostgresDumpApp = "/usr/bin/pg_dump"

// PostgresDumpallApp points to the pg_dumpall binary location
var PostgresDumpallApp = "/usr/bin/pg_dumpall"

// PostgresRestoreApp points to the pg_restore binary location
var PostgresRestoreApp = "/usr/bin/pg_restore"

// PostgresTermApp points to the psql binary location
var PostgresTermApp = "/usr/bin/psql"

var terminateQuery = `SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();`

var dropQuery = `DROP DATABASE "%s";`

var createQuery = `CREATE DATABASE "%s" OWNER "%s";`

var maintenanceDatabase = "postgres"

func (p *PostgresConfig) newBaseArgs() []string {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
	}

	if p.Database != "" {
		args = append(args, "-d", p.Database)
	}

	options := strings.Fields(p.Options)

	// add extra options
	if len(options) > 0 {
		args = append(args, options...)
	}

	return args
}

func (p *PostgresConfig) newPostgresCmd() *CmdConfig {
	var env []string

	if p.Password != "" {
		env = append(env, "PGPASSWORD="+p.Password)
	}

	return &CmdConfig{
		Env: env,
	}
}

// Backup generates a dump of the database and returns the path where is stored
func (p *PostgresConfig) Backup() (string, error) {
	var prefix string
	if p.NameAsPrefix {
		prefix = p.Database
	} else if p.NamePrefix != "" {
		prefix = p.NamePrefix
	} else {
		prefix = "postgres-backup"
	}
	filepath := generateFilename(p.SaveDir, prefix)
	args := p.newBaseArgs()

	var appPath string
	if p.Database != "" {
		appPath = PostgresDumpApp
	} else {
		appPath = PostgresDumpallApp
	}

	// only allow custom format when dumping a single database
	if p.Custom && p.Database != "" {
		filepath += ".dump"
		args = append(args, "-f", filepath)
		args = append(args, "-Fc")
	} else if !p.Compress {
		filepath += ".sql"
		args = append(args, "-f", filepath)
	} else {
		filepath += ".sql.gz"
	}

	app := p.newPostgresCmd()

	if p.Compress && !p.Custom {
		f, err := os.Create(filepath)
		if err != nil {
			return "", fmt.Errorf("cannot create file: %v", err)
		}

		defer f.Close()

		writer := gzip.NewWriter(f)
		defer writer.Close()

		app.OutputFile = writer
	}

	if err := app.CmdRun(appPath, args...); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", appPath, err)
	}

	return filepath, nil
}

// Restore takes a database dump and restores it
func (p *PostgresConfig) Restore(filepath string) error {
	args := p.newBaseArgs()
	var appPath string

	// only allow custom format when restoring a single database
	if p.Custom && p.Database != "" {
		args = append(args, filepath)
		appPath = PostgresRestoreApp
	} else {
		appPath = PostgresTermApp
	}

	app := p.newPostgresCmd()

	if !p.Custom {
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
		defer f.Close()
	}

	if p.Drop {
		log.Info("Recreating database %s", p.Database)
		if err := p.recreate(); err != nil {
			return fmt.Errorf("couldn't recreate database, %v", err)
		}
	}

	if err := app.CmdRun(appPath, args...); err != nil {
		serr, ok := err.(*exec.ExitError)

		if ok && p.IgnoreExitCode {
			log.Info("Ignored exit code of restore process: %v", serr)
		} else {
			return fmt.Errorf("couldn't execute %s, %v", appPath, err)
		}
	}

	return nil
}

func (p *PostgresConfig) recreate() error {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
		maintenanceDatabase,
	}

	app := p.newPostgresCmd()

	terminate := append(args, "-c", fmt.Sprintf(terminateQuery, p.Database))
	if err := app.CmdRun(PostgresTermApp, terminate...); err != nil {
		return fmt.Errorf("psql error on terminate, %v", err)
	}

	remove := append(args, "-c", fmt.Sprintf(dropQuery, p.Database))
	if err := app.CmdRun(PostgresTermApp, remove...); err != nil {
		return fmt.Errorf("psql error on drop, %v", err)
	}

	var owner string
	if p.Owner != "" {
		owner = p.Owner
	} else {
		owner = p.User
	}

	create := append(args, "-c", fmt.Sprintf(createQuery, p.Database, owner))
	if err := app.CmdRun(PostgresTermApp, create...); err != nil {
		return fmt.Errorf("psql error on create, %v", err)
	}

	return nil
}
