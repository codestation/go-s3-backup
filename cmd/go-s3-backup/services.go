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
	"github.com/urfave/cli"
	"megpoid.xyz/go/go-s3-backup/services"
)

var gogsFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "gogs-config",
		Usage:  "gogs config path",
		EnvVar: "GOGS_CONFIG",
	},
	cli.StringFlag{
		Name:   "gogs-data",
		Usage:  "gogs data path",
		Value:  "/data",
		EnvVar: "GOGS_DATA",
	},
}

var databaseFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "host",
		Usage:  "database host",
		EnvVar: "DATABASE_HOST",
	},
	cli.StringFlag{
		Name:   "port",
		Usage:  "database port",
		EnvVar: "DATABASE_PORT",
	},
	cli.StringFlag{
		Name:   "database",
		Usage:  "database name",
		EnvVar: "DATABASE_NAME",
	},
	cli.StringFlag{
		Name:   "user",
		Usage:  "database user",
		EnvVar: "DATABASE_USER",
	},
	cli.StringFlag{
		Name:   "password",
		Usage:  "database password",
		EnvVar: "DATABASE_PASSWORD",
	},
	cli.StringSliceFlag{
		Name:   "options",
		Usage:  "extra options to pass to database service",
		EnvVar: "DATABASE_OPTIONS",
	},
	cli.BoolFlag{
		Name:   "compress",
		Usage:  "compress sql with gzip",
		EnvVar: "DATABASE_COMPRESS",
	},
}

var postgresFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "postgres-custom",
		Usage:  "use custom format (always compressed)",
		EnvVar: "POSTGRES_CUSTOM_FORMAT",
	},
}

var tarballFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "tarball-path",
		Usage:  "path to backup/restore",
		EnvVar: "TARBALL_PATH_SOURCE",
	},
	cli.StringFlag{
		Name:   "tarball-name",
		Usage:  "backup file prefix",
		EnvVar: "TARBALL_NAME_PREFIX",
	},
	cli.BoolFlag{
		Name:   "tarball-compress",
		Usage:  "compress tarball with gzip",
		EnvVar: "TARBALL_COMPRESS",
	},
}

func newGogsConfig(c *cli.Context) *services.Gogs {
	c = c.Parent()

	return &services.Gogs{
		ConfigPath: c.String("gogs-config"),
		DataPath:   c.String("gogs-data"),
	}
}

func newMysqlConfig(c *cli.Context) *services.MySQL {
	c = c.Parent()

	return &services.MySQL{
		Host:     c.String("host"),
		Port:     c.String("port"),
		User:     c.String("user"),
		Password: c.String("password"),
		Database: c.String("database"),
		Options:  c.StringSlice("options"),
		Compress: c.Bool("compress"),
	}
}

func newPostgresConfig(c *cli.Context) *services.Postgres {
	c = c.Parent()

	return &services.Postgres{
		Host:     c.String("host"),
		Port:     c.String("port"),
		User:     c.String("user"),
		Password: c.String("password"),
		Database: c.String("database"),
		Options:  c.StringSlice("options"),
		Compress: c.Bool("compress"),
		Custom:   c.Bool("postgres-custom"),
	}
}

func newTarballConfig(c *cli.Context) *services.Tarball {
	c = c.Parent()

	return &services.Tarball{
		Path:     c.String("tarball-path"),
		Name:     c.String("tarball-name"),
		Compress: c.Bool("tarball-compress"),
	}
}

func gogsCmd(parent string) cli.Command {
	name := "gogs"
	return cli.Command{
		Name:  name,
		Usage: "connect to gogs service",
		Flags: gogsFlags,
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func postgresCmd(parent string) cli.Command {
	name := "postgres"
	return cli.Command{
		Name:  name,
		Usage: "connect to postgres service",
		Flags: append(databaseFlags, postgresFlags...),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func mysqlCmd(parent string) cli.Command {
	name := "mysql"
	return cli.Command{
		Name:  name,
		Usage: "connect to mysql service",
		Flags: databaseFlags,
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func tarballCmd(parent string) cli.Command {
	name := "tarball"
	return cli.Command{
		Name:  name,
		Usage: "connect to tarball service",
		Flags: tarballFlags,
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}
