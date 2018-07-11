package main

import (
	"github.com/urfave/cli"
	"megpoid.xyz/go/go-s3-backup/services"
)

var GogsFlags = []cli.Flag{
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

var DatabaseFlags = []cli.Flag{
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

var PostgresFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "postgres-custom",
		Usage:  "use custom format",
		EnvVar: "POSTGRES_CUSTOM_FORMAT",
	},
}

var TarballFlags = []cli.Flag{
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

func NewGogsConfig(c *cli.Context) *services.Gogs {
	return &services.Gogs{
		ConfigPath: c.String("gogs-config"),
		DataPath:   c.String("gogs-data"),
	}
}

func NewMysqlConfig(c *cli.Context) *services.MySQL {
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

func NewPostgresConfig(c *cli.Context) *services.Postgres {
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

func NewTarballConfig(c *cli.Context) *services.Tarball {
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
		Flags: GogsFlags,
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
		Flags: append(DatabaseFlags, PostgresFlags...),
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
		Flags: DatabaseFlags,
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
		Flags: TarballFlags,
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}
