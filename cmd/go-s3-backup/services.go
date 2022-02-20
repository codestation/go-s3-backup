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

package main

import (
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"megpoid.xyz/go/go-s3-backup/services"
)

var giteaFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "gitea-config",
		Usage:  "gitea config path",
		EnvVar: "GOGS_CONFIG",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "gitea-data",
		Usage:  "gitea data path",
		Value:  "/data",
		EnvVar: "GOGS_DATA",
	}),
}

var databaseFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-host",
		Usage:  "database host",
		EnvVar: "DATABASE_HOST",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-port",
		Usage:  "database port",
		EnvVar: "DATABASE_PORT",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-name",
		Usage:  "database name",
		EnvVar: "DATABASE_NAME",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-filename-prefix",
		Usage:  "database filename prefix",
		EnvVar: "DATABASE_FILENAME_PREFIX",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "database-name-as-prefix",
		Usage:  "database name as prefix",
		EnvVar: "DATABASE_NAME_AS_PREFIX",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-user",
		Usage:  "database user",
		EnvVar: "DATABASE_USER",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-password",
		Usage:  "database password",
		EnvVar: "DATABASE_PASSWORD",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-password-file",
		Usage:  "database password file",
		EnvVar: "DATABASE_PASSWORD_FILE",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-options",
		Usage:  "extra options to pass to database service",
		EnvVar: "DATABASE_OPTIONS",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "database-compress",
		Usage:  "compress sql with gzip",
		EnvVar: "DATABASE_COMPRESS",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "database-ignore-exit-code",
		Usage:  "ignore restore process exit code",
		EnvVar: "DATABASE_IGNORE_EXIT_CODE",
	}),
}

var postgresFlags = []cli.Flag{
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "postgres-custom",
		Usage:  "use custom format (always compressed), ignored when database name is not set",
		EnvVar: "POSTGRES_CUSTOM_FORMAT",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "postgres-drop",
		Usage:  "drop database before restoring it",
		EnvVar: "POSTGRES_DROP",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "postgres-owner",
		Usage:  "change owner on database restore",
		EnvVar: "POSTGRES_OWNER",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:   "postgres-exclude-databases",
		Usage:  "make backup of all databases except the ones listed",
		EnvVar: "POSTGRES_EXCLUDE_DATABASES",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "postgres-backup-per-user",
		Usage:  "make backups for all databases separated per user",
		EnvVar: "POSTGRES_BACKUP_PER_USER",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:   "postgres-backup-users",
		Usage:  "make backups for databases matching these users",
		EnvVar: "POSTGRES_BACKUP_USERS",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:   "postgres-backup-exclude-users",
		Usage:  "make backups for databases excluding these users",
		EnvVar: "POSTGRES_BACKUP_EXCLUDE_USERS",
		Value:  &cli.StringSlice{"postgres"},
	}),
}

var tarballFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "tarball-path",
		Usage:  "path to backup/restore",
		EnvVar: "TARBALL_PATH_SOURCE",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "tarball-name",
		Usage:  "backup file prefix",
		EnvVar: "TARBALL_NAME_PREFIX",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "tarball-compress",
		Usage:  "compress tarball with gzip",
		EnvVar: "TARBALL_COMPRESS",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "tarball-path-prefix",
		Usage:  "backup path prefix",
		EnvVar: "TARBALL_PATH_PREFIX",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "tarball-backup-per-dir",
		Usage:  "backup each folder individually",
		EnvVar: "TARBALL_BACKUP_PER_DIR",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:   "tarball-backup-dirs",
		Usage:  "backup each folder individually",
		EnvVar: "TARBALL_BACKUP_DIRS",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:   "tarball-backup-exclude-dirs",
		Usage:  "make backups for directories excluding these dirs",
		EnvVar: "TARBALL_BACKUP_EXCLUDE_DIRS",
	}),
}

func newGogsConfig(c *cli.Context) *services.GiteaConfig {
	c = c.Parent()

	return &services.GiteaConfig{
		ConfigPath: c.String("gitea-config"),
		DataPath:   c.String("gitea-data"),
		SaveDir:    c.GlobalString("savedir"),
	}
}

func newMysqlConfig(c *cli.Context) *services.MySQLConfig {
	c = c.Parent()

	return &services.MySQLConfig{
		Host:           c.String("database-host"),
		Port:           c.String("database-port"),
		User:           c.String("database-user"),
		Password:       fileOrString(c, "database-password"),
		Database:       c.String("database-name"),
		NamePrefix:     c.String("database-filename-prefix"),
		NameAsPrefix:   c.Bool("database-name-as-prefix"),
		Options:        c.String("database-options"),
		Compress:       c.Bool("database-compress"),
		SaveDir:        c.GlobalString("savedir"),
		IgnoreExitCode: c.Bool("database-ignore-exit-code"),
	}
}

func newPostgresConfig(c *cli.Context) *services.PostgresConfig {
	c = c.Parent()

	return &services.PostgresConfig{
		Host:             c.String("database-host"),
		Port:             c.String("database-port"),
		User:             c.String("database-user"),
		Password:         fileOrString(c, "database-password"),
		Database:         c.String("database-name"),
		NamePrefix:       c.String("database-filename-prefix"),
		NameAsPrefix:     c.Bool("database-name-as-prefix"),
		Options:          c.String("database-options"),
		Compress:         c.Bool("database-compress"),
		Custom:           c.Bool("postgres-custom"),
		SaveDir:          c.GlobalString("savedir"),
		IgnoreExitCode:   c.Bool("database-ignore-exit-code"),
		Drop:             c.Bool("postgres-drop"),
		Owner:            c.String("postgres-owner"),
		ExcludeDatabases: c.StringSlice("postgres-exclude-databases"),
		BackupPerUser:    c.Bool("postgres-backup-per-user"),
		BackupUsers:      c.StringSlice("postgres-backup-users"),
		ExcludeUsers:     c.StringSlice("postgres-backup-exclude-users"),
	}
}

func newTarballConfig(c *cli.Context) *services.TarballConfig {
	c = c.Parent()

	return &services.TarballConfig{
		Name:         c.String("tarball-name"),
		Path:         c.String("tarball-path"),
		Compress:     c.Bool("tarball-compress"),
		SaveDir:      c.GlobalString("savedir"),
		Prefix:       c.String("tarball-path-prefix"),
		BackupPerDir: c.Bool("tarball-backup-per-dir"),
		BackupDirs:   c.StringSlice("tarball-backup-dirs"),
		ExcludeDirs:  c.StringSlice("tarball-backup-exclude-dirs"),
	}
}

func newConsulConfig(c *cli.Context) *services.ConsulConfig {
	c = c.Parent()

	return &services.ConsulConfig{
		SaveDir: c.GlobalString("savedir"),
	}
}

func giteaCmd(parent string) cli.Command {
	name := "gitea"
	return cli.Command{
		Name:   name,
		Usage:  "connect to gitea service",
		Flags:  giteaFlags,
		Before: applyConfigValues(giteaFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func postgresCmd(parent string) cli.Command {
	name := "postgres"
	flags := append(databaseFlags, postgresFlags...)
	return cli.Command{
		Name:   name,
		Usage:  "connect to postgres service",
		Flags:  flags,
		Before: applyConfigValues(flags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func mysqlCmd(parent string) cli.Command {
	name := "mysql"
	return cli.Command{
		Name:   name,
		Usage:  "connect to mysql service",
		Flags:  databaseFlags,
		Before: applyConfigValues(databaseFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func tarballCmd(parent string) cli.Command {
	name := "tarball"
	return cli.Command{
		Name:   name,
		Usage:  "connect to tarball service",
		Flags:  tarballFlags,
		Before: applyConfigValues(tarballFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func consulCmd(parent string) cli.Command {
	name := "consul"
	return cli.Command{
		Name:  name,
		Usage: "connect to consul service",
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}
