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

package commands

import (
	"github.com/spf13/viper"
	"go.megpoid.dev/go-s3-backup/stores"
)

func newS3Config() *stores.S3Config {
	return &stores.S3Config{
		// S3 config
		Endpoint:        viper.GetString("s3-endpoint"),
		Region:          viper.GetString("s3-region"),
		Bucket:          viper.GetString("s3-bucket"),
		Prefix:          viper.GetString("s3-prefix"),
		ForcePathStyle:  viper.GetBool("s3-force-path-style"),
		KeepAfterUpload: viper.GetBool("s3-keep-file"),
		// default config
		SaveDir: viper.GetString("save-dir"),
	}
}

func newFilesystemConfig() *stores.FilesystemConfig {
	return &stores.FilesystemConfig{
		// default config
		SaveDir: viper.GetString("save-dir"),
	}
}
