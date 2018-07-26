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
	"megpoid.xyz/go/go-s3-backup/stores"
)

var s3Flags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "s3-endpoint",
		Usage:  "s3 endpoint",
		EnvVar: "S3_ENDPOINT",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "s3-region",
		Usage:  "s3 region",
		EnvVar: "S3_REGION",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "s3-bucket",
		Usage:  "s3 bucket",
		EnvVar: "S3_BUCKET",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "s3-prefix",
		Usage:  "s3 prefix",
		EnvVar: "S3_PREFIX",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "s3-force-path-style",
		Usage:  "s3 force path style (needed for minio)",
		EnvVar: "S3_FORCE_PATH_STYLE",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "s3-keep-file",
		Usage:  "keep local file after successful upload",
		EnvVar: "S3_KEEP_FILE",
	}),
}

func newS3Config(c *cli.Context) *stores.S3Config {
	return &stores.S3Config{
		Endpoint:        c.String("s3-endpoint"),
		Region:          c.String("s3-region"),
		Bucket:          c.String("s3-bucket"),
		Prefix:          c.String("s3-prefix"),
		ForcePathStyle:  c.Bool("s3-force-path-style"),
		KeepAfterUpload: c.Bool("s3-keep-file"),
		SaveDir:         c.GlobalString("savedir"),
	}
}

func newFilesystemConfig(c *cli.Context) *stores.FilesystemConfig {
	return &stores.FilesystemConfig{
		SaveDir: c.GlobalString("savedir"),
	}
}

func s3Cmd(command string, service string) cli.Command {
	name := "s3"
	return cli.Command{
		Name:   name,
		Usage:  "use S3Config as store",
		Flags:  s3Flags,
		Before: applyConfigValues(s3Flags),
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
