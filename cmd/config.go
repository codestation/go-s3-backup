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

package cmd

import (
	"megpoid.xyz/go/go-s3-backup/services"
	"megpoid.xyz/go/go-s3-backup/stores"

	"github.com/urfave/cli"
)

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
		Custom:   c.Bool("custom"),
	}
}

func NewS3Config(c *cli.Context) *stores.S3 {
	return &stores.S3{
		Endpoint:       c.GlobalString("endpoint"),
		Region:         c.GlobalString("region"),
		Bucket:         c.GlobalString("bucket"),
		AccessKey:      c.GlobalString("access-key"),
		ClientSecret:   c.GlobalString("client-secret"),
		ForcePathStyle: c.GlobalBool("force-path-style"),
	}
}
