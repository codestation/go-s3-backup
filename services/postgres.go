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

// Postgres has the config options for the Postgres service
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  []string
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

// Backup generates a dump of the database and returns the path where is stored
func (p *Postgres) Backup() (string, error) {
	filepath := fmt.Sprintf("%s/postgres-backup-%s", p.SaveDir, time.Now().Format("20060102150405"))
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
	}

	var app string

	if p.Database != "" {
		args = append(args, "-d", p.Database)
		app = PostgresDumpApp
	} else {
		app = PostgresDumpallApp
	}

	env := os.Environ()

	if p.Password != "" {
		env = append(env, "PGPASSWORD="+p.Password)
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

	// add extra options
	if len(p.Options) > 0 {
		args = append(args, p.Options...)
	}

	cmd := exec.Command(app, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if p.Compress {
		if err := compressAppOutput(cmd, filepath); err != nil {
			return "", fmt.Errorf("couldn't pipe app output, %v", err)
		}
	} else {
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("couldn't execute %s, %v", PostgresDumpApp, err)
		}
	}

	return filepath, nil
}

// Restore takes a database dump and restores it
func (p *Postgres) Restore(filepath string) error {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
	}

	if p.Database != "" {
		args = append(args, "-d", p.Database)
	}

	env := os.Environ()

	if p.Password != "" {
		env = append(env, "PGPASSWORD="+p.Password)
	}

	// add extra options
	if len(p.Options) > 0 {
		args = append(args, p.Options...)
	}

	var App string

	// only allow custom format when restoring a single database
	if p.Custom && p.Database != "" {
		args = append(args, filepath)
		App = PostgresRestoreApp
	} else {
		App = PostgresTermApp
	}

	cmd := exec.Command(App, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if !p.Custom {
		if err := readFileToInput(cmd, filepath); err != nil {
			return fmt.Errorf("couldn't pipe file to input, %v", err)
		}
	} else {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("couldn't execute pg_restore, %v", err)
		}
	}

	return nil
}
