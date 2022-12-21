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

// PostgresConfig has the config options for the PostgresConfig service
type PostgresConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	Database         string
	NamePrefix       string
	NameAsPrefix     bool
	Options          string
	Compress         bool
	Custom           bool
	SaveDir          string
	IgnoreExitCode   bool
	Drop             bool
	Owner            string
	ExcludeDatabases []string
	BackupPerUser    bool
	BackupUsers      []string
	ExcludeUsers     []string
	BackupPerSchema  bool
	BackupSchemas    []string
	ExcludeSchemas   []string
	Version          string
}

// PostgresBinaryPath points to the location where the postgres binaries are located
var PostgresBinaryPath = "/usr/bin"

var terminateQuery = `SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity
WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();`

var dropQuery = `DROP DATABASE "%s";`

var createQuery = `CREATE DATABASE "%s" OWNER "%s";`

var postgresListDatabasesQuery = `COPY(SELECT datname FROM pg_database JOIN pg_authid ON pg_database.datdba = pg_authid.oid
WHERE rolname = '%s' ORDER BY datname) TO STDOUT`

var listUsersQuery = `COPY(SELECT usename FROM pg_catalog.pg_user ORDER BY usename) TO STDOUT;`

var listSchemasQuery = `COPY(SELECT nspname FROM pg_catalog.pg_namespace ORDER BY nspname) TO STDOUT;`

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

func (p *PostgresConfig) Backup() (*BackupResults, error) {
	switch {
	case p.BackupPerUser:
		return p.backupPerUser()
	case p.BackupPerSchema:
		return p.backupPerSchema()
	default:
		namePrefix := p.getNamePrefix()
		filepath, err := p.backupDatabase("", namePrefix)
		if err != nil {
			return nil, err
		}

		return &BackupResults{Entries: []BackupResult{{
			NamePrefix: namePrefix,
			Path:       filepath,
		}}}, nil
	}
}

// Backup generates a dump of the database and returns the path where is stored
func (p *PostgresConfig) backupPerUser() (*BackupResults, error) {
	users, err := p.listUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to list users, %w", err)
	}

	var resultList []BackupResult

	for _, user := range users {
		found := false
		for _, u := range p.ExcludeUsers {
			if u == user {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if len(p.BackupUsers) > 0 {
			found := false
			for _, u := range p.BackupUsers {
				if u == user {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		databases, err := p.listDatabases(user)
		if err != nil {
			return nil, fmt.Errorf("failed to list databases for user %s, %w", user, err)
		}

		for _, database := range databases {

			if len(p.ExcludeDatabases) > 0 {
				found := false
				for _, exclude := range p.ExcludeDatabases {
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

			p.Database = database
			namePrefix := p.getNamePrefix()
			filepath, err := p.backupDatabase(user, namePrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to backup database %s, %w", database, err)
			}
			resultEntry := BackupResult{DirPrefix: user, NamePrefix: namePrefix, Path: filepath}
			resultList = append(resultList, resultEntry)
		}
	}

	result := &BackupResults{resultList}
	return result, nil
}

// Backup generates a dump of the database and returns the path where is stored
func (p *PostgresConfig) backupPerSchema() (*BackupResults, error) {
	schemas, err := p.listSchemas(p.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas, %w", err)
	}

	var resultList []BackupResult

	for _, schema := range schemas {
		found := false
		for _, u := range p.ExcludeSchemas {
			if u == schema {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if len(p.BackupSchemas) > 0 {
			found := false
			for _, u := range p.BackupSchemas {
				if u == schema {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		namePrefix := p.getNamePrefix() + "_" + schema
		baseDir := path.Join(p.Database, schema)
		filepath, err := p.backupDatabase(baseDir, namePrefix, schema)
		if err != nil {
			return nil, fmt.Errorf("failed to backup database schema %s, %w", schema, err)
		}
		resultEntry := BackupResult{DirPrefix: baseDir, NamePrefix: namePrefix, Path: filepath}
		resultList = append(resultList, resultEntry)
	}

	result := &BackupResults{resultList}
	return result, nil
}

func (p *PostgresConfig) getNamePrefix() string {
	var prefix string
	if p.NameAsPrefix && p.Database != "" {
		prefix = p.Database
	} else if p.NamePrefix != "" {
		prefix = p.NamePrefix
	} else {
		prefix = "postgres-backup"
	}

	return prefix
}

// Backup generates a dump of the database and returns the path where is stored
func (p *PostgresConfig) backupDatabase(basedir, namePrefix string, schemas ...string) (string, error) {
	savePath := path.Join(p.SaveDir, basedir)
	filepath := generateFilename(savePath, namePrefix)
	args := p.newBaseArgs()

	var appPath string
	if p.Database != "" {
		appPath = path.Join(PostgresBinaryPath, "pg_dump")
		for _, schema := range schemas {
			args = append(args, "--schema="+schema)
		}
	} else {
		appPath = path.Join(PostgresBinaryPath, "pg_dumpall")
		// no custom format for pg_dumpall
		p.Custom = false
		if len(p.ExcludeDatabases) > 0 {
			for _, exclude := range p.ExcludeDatabases {
				args = append(args, "--exclude-database="+exclude)
			}
		}
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

	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", err
	}

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
		appPath = path.Join(PostgresBinaryPath, "pg_restore")
	} else {
		appPath = path.Join(PostgresBinaryPath, "psql")
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
	psqlApp := path.Join(PostgresBinaryPath, "psql")

	terminate := append(args, "-c", fmt.Sprintf(terminateQuery, p.Database))
	if err := app.CmdRun(psqlApp, terminate...); err != nil {
		return fmt.Errorf("psql error on terminate, %v", err)
	}

	remove := append(args, "-c", fmt.Sprintf(dropQuery, p.Database))
	if err := app.CmdRun(psqlApp, remove...); err != nil {
		return fmt.Errorf("psql error on drop, %v", err)
	}

	var owner string
	if p.Owner != "" {
		owner = p.Owner
	} else {
		owner = p.User
	}

	create := append(args, "-c", fmt.Sprintf(createQuery, p.Database, owner))
	if err := app.CmdRun(psqlApp, create...); err != nil {
		return fmt.Errorf("psql error on create, %v", err)
	}

	return nil
}

func (p *PostgresConfig) listDatabases(user string) ([]string, error) {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
		maintenanceDatabase,
	}

	app := p.newPostgresCmd()
	psqlApp := path.Join(PostgresBinaryPath, "psql")

	var b bytes.Buffer
	outputWriter := bufio.NewWriter(&b)
	app.OutputFile = outputWriter

	listDatabases := append(args, "-c", fmt.Sprintf(postgresListDatabasesQuery, user))
	if err := app.CmdRun(psqlApp, listDatabases...); err != nil {
		return nil, fmt.Errorf("psql error on database list, %w", err)
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

func (p *PostgresConfig) listUsers() ([]string, error) {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
		maintenanceDatabase,
	}

	app := p.newPostgresCmd()
	psqlApp := path.Join(PostgresBinaryPath, "psql")

	var b bytes.Buffer
	outputWriter := bufio.NewWriter(&b)
	app.OutputFile = outputWriter

	listUsers := append(args, "-c", listUsersQuery)
	if err := app.CmdRun(psqlApp, listUsers...); err != nil {
		return nil, fmt.Errorf("psql error on user list, %w", err)
	}

	scanner := bufio.NewScanner(&b)
	var users []string
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse psql output, %w", err)
	}

	return users, nil
}

func (p *PostgresConfig) listSchemas(database string) ([]string, error) {
	args := []string{
		"-h", p.Host,
		"-p", p.Port,
		"-U", p.User,
		database,
	}

	app := p.newPostgresCmd()
	psqlApp := path.Join(PostgresBinaryPath, "psql")

	var b bytes.Buffer
	outputWriter := bufio.NewWriter(&b)
	app.OutputFile = outputWriter

	listSchemas := append(args, "-c", listSchemasQuery)
	if err := app.CmdRun(psqlApp, listSchemas...); err != nil {
		return nil, fmt.Errorf("psql error on schema list, %w", err)
	}

	scanner := bufio.NewScanner(&b)
	var users []string
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse psql output, %w", err)
	}

	return users, nil
}
