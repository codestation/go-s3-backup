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
	"strings"
)

// PostgresConfig has the config options for the PostgresConfig service
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string
	Compress bool
	Custom   bool
	SaveDir  string
}

// PostgresDumpApp points to the pg_dump binary location
var PostgresDumpApp = "/usr/bin/pg_dump"

// PostgresDumpallApp points to the pg_dumpall binary location
var PostgresDumpallApp = "/usr/bin/pg_dumpall"

// PostgresRestoreApp points to the pg_restore binary location
var PostgresRestoreApp = "/usr/bin/pg_restore"

// PostgresTermApp points to the psql binary location
var PostgresTermApp = "/usr/bin/psql"

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
	filepath := generateFilename(p.SaveDir, "postgres-backup")
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

	if p.Compress {
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
		return "", fmt.Errorf("couldn't execute %s, %v", PostgresDumpApp, err)
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

	if err := app.CmdRun(appPath, args...); err != nil {
		return fmt.Errorf("couldn't execute %s, %v", appPath, err)
	}

	return nil
}
