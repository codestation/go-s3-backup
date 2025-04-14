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
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackupRestore(t *testing.T) {
	r := require.New(t)
	tmp, err := os.MkdirTemp("", "archiver")
	r.NoError(err, "failed to create temp directory")

	defer os.RemoveAll(tmp)

	backupDir := path.Join(tmp, "backup")
	err = os.Mkdir(backupDir, 0o755)
	r.NoError(err, "failed to create backup directory")

	filepath := path.Join(backupDir, "test.txt")
	expected := []byte("test")
	err = os.WriteFile(filepath, expected, 0o777)
	r.NoError(err, "failed to create backup file")

	tar := TarballConfig{
		Path:     backupDir,
		Name:     "test",
		Compress: true,
		SaveDir:  tmp,
	}

	results, err := tar.Backup()
	r.NoError(err, "failed to create backup tarball")

	err = os.RemoveAll(backupDir)
	r.NoError(err, "failed to remove backup directory")

	err = os.Mkdir(backupDir, 0o755)
	r.NoError(err, "failed to create backup directory")

	for _, result := range results.Entries {
		err = tar.Restore(result.Path)
		r.NoError(err, "failed to restore backup dir")

		actual, err := os.ReadFile(filepath)
		r.NoError(err, "failed to read restored file")
		r.Equal(expected, actual, "backup contents mismatch")
	}
}
