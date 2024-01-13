package stores

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreRestore(t *testing.T) {
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

	fs := FilesystemConfig{
		SaveDir: tmp,
	}

	err = fs.Store(filepath, "", "test.txt")
	r.NoError(err, "failed to store file")
}
