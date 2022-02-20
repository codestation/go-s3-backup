package services

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackupRestore(t *testing.T) {
	r := require.New(t)
	tmp, err := ioutil.TempDir("", "archiver")
	r.NoError(err, "failed to create temp directory")

	defer os.RemoveAll(tmp)

	backupDir := path.Join(tmp, "backup")
	err = os.Mkdir(backupDir, 0755)
	r.NoError(err, "failed to create backup directory")

	filepath := path.Join(backupDir, "test.txt")
	expected := []byte("test")
	err = ioutil.WriteFile(filepath, expected, 0777)
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
		for _, tarball := range result.Filenames {
			err = tar.Restore(tarball)
			r.NoError(err, "failed to restore backup dir")

			actual, err := ioutil.ReadFile(filepath)
			r.NoError(err, "failed to read restored file")
			r.Equal(expected, actual, "backup contents mismatch")
		}
	}
}
