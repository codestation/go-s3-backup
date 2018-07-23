package stores

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreRestore(t *testing.T) {
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

	fs := FilesystemConfig{
		SaveDir: tmp,
	}

	err = fs.Store(filepath, "test.txt")
	r.NoError(err, "failed to store file")
}
