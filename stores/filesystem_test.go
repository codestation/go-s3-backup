package stores


import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreRestore(t *testing.T) {
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


	fs := FilesystemConfig{
		SaveDir: tmp,
	}

	err = fs.Store(filepath, "test.txt")
	assert.Nil(t, err, "failed to store file")
}
