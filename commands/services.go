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
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"go.megpoid.dev/go-s3-backup/services"
)

// https://github.com/spf13/viper/issues/380#issuecomment-1916465489
func getStringSlice(key string) []string {
	var h []string
	if err := viper.UnmarshalKey(key, &h); err != nil {
		slog.Error("Failed to unmarshal key", "key", key, "error", err)
		return []string{}
	}
	return h
}

func newMysqlConfig() *services.MySQLConfig {
	return &services.MySQLConfig{
		// database config
		Host:           viper.GetString("database-host"),
		Port:           viper.GetString("database-port"),
		Database:       viper.GetString("database-name"),
		User:           viper.GetString("database-user"),
		Password:       fileOrString("database-password"),
		NamePrefix:     viper.GetString("database-filename-prefix"),
		NameAsPrefix:   viper.GetBool("database-name-as-prefix"),
		Options:        viper.GetString("database-options"),
		Compress:       viper.GetBool("database-compress"),
		IgnoreExitCode: viper.GetBool("database-ignore-exit-code"),
		// mysql config
		SplitDatabases:   viper.GetBool("mysql-split-databases"),
		ExcludeDatabases: getStringSlice("mysql-exclude-databases"),
		// default config
		SaveDir: viper.GetString("save-dir"),
	}
}

func newPostgresConfig() *services.PostgresConfig {
	services.PostgresBinaryPath = viper.GetString("postgres-binary-path")
	if services.PostgresBinaryPath == "" {
		services.PostgresBinaryPath = fmt.Sprintf("/usr/libexec/postgresql%s", viper.GetString("postgres-version"))
	}

	return &services.PostgresConfig{
		// database config
		Host:           viper.GetString("database-host"),
		Port:           viper.GetString("database-port"),
		Database:       viper.GetString("database-name"),
		User:           viper.GetString("database-user"),
		Password:       fileOrString("database-password"),
		NamePrefix:     viper.GetString("database-filename-prefix"),
		NameAsPrefix:   viper.GetBool("database-name-as-prefix"),
		Options:        viper.GetString("database-options"),
		Compress:       viper.GetBool("database-compress"),
		IgnoreExitCode: viper.GetBool("database-ignore-exit-code"),
		// postgres config
		BinaryPath:       services.PostgresBinaryPath,
		Version:          viper.GetString("postgres-version"),
		Custom:           viper.GetBool("postgres-custom-format"),
		Drop:             viper.GetBool("postgres-drop"),
		Owner:            viper.GetString("postgres-owner"),
		ExcludeDatabases: getStringSlice("postgres-exclude-databases"),
		BackupPerUser:    viper.GetBool("postgres-backup-per-user"),
		BackupUsers:      getStringSlice("postgres-backup-users"),
		ExcludeUsers:     getStringSlice("postgres-backup-exclude-users"),
		BackupPerSchema:  viper.GetBool("postgres-backup-per-schema"),
		BackupSchemas:    getStringSlice("postgres-backup-schemas"),
		ExcludeSchemas:   getStringSlice("postgres-backup-exclude-schemas"),
		// default config
		SaveDir: viper.GetString("save-dir"),
	}
}

func newTarballConfig() *services.TarballConfig {
	return &services.TarballConfig{
		// tarball config
		Name:         viper.GetString("tarball-name-prefix"),
		Path:         viper.GetString("tarball-path-source"),
		Compress:     viper.GetBool("tarball-compress"),
		Prefix:       viper.GetString("tarball-path-prefix"),
		BackupPerDir: viper.GetBool("tarball-backup-per-dir"),
		BackupDirs:   getStringSlice("tarball-backup-dirs"),
		ExcludeDirs:  getStringSlice("tarball-backup-exclude-dirs"),
		// default config
		SaveDir: viper.GetString("save-dir"),
	}
}
