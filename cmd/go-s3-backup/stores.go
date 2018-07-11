package main

import (
	"github.com/urfave/cli"
	"megpoid.xyz/go/go-s3-backup/stores"
)

var S3Flags = []cli.Flag{
	cli.StringFlag{
		Name:   "s3-endpoint",
		Usage:  "s3 endpoint",
		EnvVar: "S3_ENDPOINT",
	},
	cli.StringFlag{
		Name:   "s3-region",
		Usage:  "s3 region",
		EnvVar: "S3_REGION",
	},
	cli.StringFlag{
		Name:   "s3-bucket",
		Usage:  "s3 bucket",
		EnvVar: "S3_BUCKET",
	},
	cli.StringFlag{
		Name:   "s3-prefix",
		Usage:  "s3 prefix",
		EnvVar: "S3_PREFIX",
	},
	cli.StringFlag{
		Name:   "s3-access-key",
		Usage:  "s3 access key",
		EnvVar: "S3_ACCESS_KEY",
	},
	cli.StringFlag{
		Name:   "s3-client secret",
		Usage:  "s3 client secret",
		EnvVar: "S3_CLIENT_SECRET",
	},
	cli.BoolFlag{
		Name:   "s3-force-path-style",
		Usage:  "s3 force path style (needed for minio)",
		EnvVar: "S3_FORCE_PATH_STYLE",
	},
	cli.BoolFlag{
		Name:   "s3-remove-after",
		Usage:  "remove file after successful upload",
		EnvVar: "S3_REMOVE_AFTER_UPLOAD",
	},
}

func NewS3Config(c *cli.Context) *stores.S3 {
	return &stores.S3{
		Endpoint:          c.GlobalString("s3-endpoint"),
		Region:            c.GlobalString("s3-region"),
		Bucket:            c.GlobalString("s3-bucket"),
		AccessKey:         c.GlobalString("s3-access-key"),
		ClientSecret:      c.GlobalString("s3-client-secret"),
		Prefix:            c.GlobalString("s3-prefix"),
		ForcePathStyle:    c.GlobalBool("s3-force-path-style"),
		RemoveAfterUpload: c.GlobalBool("s3-remove-after"),
	}
}

func NewFilesystemConfig(c *cli.Context) *stores.Filesystem {
	return &stores.Filesystem{
		SaveDir: c.GlobalString("savedir"),
	}
}

func s3Cmd(command string, service string) cli.Command {
	name := "s3"
	return cli.Command{
		Name:  name,
		Usage: "use S3 as store",
		Flags: S3Flags,
		Action: func(c *cli.Context) error {
			return runTask(c, command, service, name)
		},
	}
}

func filesystemCmd(command string, service string) cli.Command {
	name := "filesystem"
	return cli.Command{
		Name:  name,
		Usage: "use the filesystem as store",
		Action: func(c *cli.Context) error {
			return runTask(c, command, service, name)
		},
	}
}
