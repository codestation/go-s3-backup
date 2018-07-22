package services

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackupRestore(t *testing.T) {
	tmp, err := ioutil.TempDir("", "archiver")
	assert.Nil(t, err, "failed to create temp directory")

	defer os.RemoveAll(tmp)

	backupDir := path.Join(tmp, "backup")
	err = os.Mkdir(backupDir, 0755)
	assert.Nil(t, err, "failed to create backup directory")

	filepath := path.Join(backupDir, "test.txt")
	expected := []byte("test")
	err = ioutil.WriteFile(filepath, expected, 0777)
	assert.Nil(t, err, "failed to create backup file")

	tar := TarballConfig{
		Path:     backupDir,
		Name:     "test",
		Compress: true,
		SaveDir:  tmp,
	}

	tarball, err := tar.Backup()
	assert.Nil(t, err, "failed to create backup tarball")

	err = tar.Restore(tarball)
	assert.Nil(t, err, "failed to restore backup dir")

	actual, err := ioutil.ReadFile(filepath)
	assert.Nil(t, err, "failed to read restored file")
	assert.Equal(t, expected, actual, "backup contents mismatch")
}
