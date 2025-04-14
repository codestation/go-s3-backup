/*
Copyright 2025 codestation

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
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
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
	BinaryPath       string
	Version          string
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
}

// PostgresBinaryPath points to the location where the postgres binaries are located
var PostgresBinaryPath = "/usr/bin"

var terminateQuery = `select pg_terminate_backend(pg_stat_activity.pid) from pg_stat_activity
where pg_stat_activity.datname = '%s' and pid <> pg_backend_pid();`

var dropQuery = `drop database "%s";`

var createQuery = `create database "%s" owner "%s";`

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
						slog.Error("Invalid pattern, skipping", "pattern", exclude)
						found = true
						break
					}
					if matched {
						slog.Info("Excluding database that match excluded pattern", "database", database, "pattern", exclude)
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

		slog.Info("Backing up schema", "name", schema)

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
	switch {
	case p.NameAsPrefix && p.Database != "":
		prefix = p.Database
	case p.NamePrefix != "":
		prefix = p.NamePrefix
	default:
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
	switch {
	case p.Custom && p.Database != "":
		filepath += ".dump"
		args = append(args, "-f", filepath)
		args = append(args, "-Fc")
	case !p.Compress:
		filepath += ".sql"
		args = append(args, "-f", filepath)
	default:
		filepath += ".sql.gz"
	}

	app := p.newPostgresCmd()

	if err := os.MkdirAll(savePath, 0o755); err != nil {
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
		slog.Info("Recreating database", "name", p.Database)
		if err := p.recreate(); err != nil {
			return fmt.Errorf("couldn't recreate database, %v", err)
		}
	}

	if err := app.CmdRun(appPath, args...); err != nil {
		serr, ok := err.(*exec.ExitError)

		if ok && p.IgnoreExitCode {
			slog.Info("Ignored exit code of restore process", "error", serr)
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

	var terminate []string
	terminate = append(terminate, args...)
	terminate = append(terminate, "-c", fmt.Sprintf(terminateQuery, p.Database))
	if err := app.CmdRun(psqlApp, terminate...); err != nil {
		return fmt.Errorf("psql error on terminate, %v", err)
	}

	var remove []string
	remove = append(remove, args...)
	remove = append(remove, "-c", fmt.Sprintf(dropQuery, p.Database))
	if err := app.CmdRun(psqlApp, remove...); err != nil {
		return fmt.Errorf("psql error on drop, %v", err)
	}

	var owner string
	if p.Owner != "" {
		owner = p.Owner
	} else {
		owner = p.User
	}

	var create []string
	create = append(create, args...)
	create = append(create, "-c", fmt.Sprintf(createQuery, p.Database, owner))
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

	var listDatabases []string
	listDatabases = append(listDatabases, args...)
	listDatabases = append(listDatabases, "-c", fmt.Sprintf(postgresListDatabasesQuery, user))
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

	var listUsers []string
	listUsers = append(listUsers, args...)
	listUsers = append(listUsers, "-c", listUsersQuery)
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

	var listSchemas []string
	listSchemas = append(listSchemas, args...)
	listSchemas = append(listSchemas, "-c", listSchemasQuery)
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
