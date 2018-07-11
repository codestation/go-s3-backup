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
		Name:   "s3-secret-key",
		Usage:  "s3 secret key",
		EnvVar: "S3_SECRET_KEY",
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
		Endpoint:          c.String("s3-endpoint"),
		Region:            c.String("s3-region"),
		Bucket:            c.String("s3-bucket"),
		AccessKey:         c.String("s3-access-key"),
		ClientSecret:      c.String("s3-secret-key"),
		Prefix:            c.String("s3-prefix"),
		ForcePathStyle:    c.Bool("s3-force-path-style"),
		RemoveAfterUpload: c.Bool("s3-remove-after"),
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
