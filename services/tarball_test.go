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
	err = os.Mkdir(backupDir, 0755)
	r.NoError(err, "failed to create backup directory")

	filepath := path.Join(backupDir, "test.txt")
	expected := []byte("test")
	err = os.WriteFile(filepath, expected, 0777)
	r.NoError(err, "failed to create backup file")

	tar := TarballConfig{
		Path:     backupDir,
		Name:     "test",
		Compress: true,
		SaveDir:  tmp,
	}

	results, err := tar.Backup()
	r.NoError(err, "failed to create backup tarball")

	for _, result := range results.Entries {
		err = tar.Restore(result.Path)
		r.NoError(err, "failed to restore backup dir")

		actual, err := os.ReadFile(filepath)
		r.NoError(err, "failed to read restored file")
		r.Equal(expected, actual, "backup contents mismatch")
	}
}
