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

package stores

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "gopkg.in/clog.v1"
)

// Filesystem has the config options for the Filesystem service
type Filesystem struct {
	SaveDir string
}

// Store moves/copies a file to another directory
func (f *Filesystem) Store(src string, filename string) error {
	dest := path.Clean(path.Join(f.SaveDir, filename))

	if src == dest {
		log.Trace("using the same path as source and destination, do nothing")
		return nil
	}

	err := os.Rename(dest, src)
	if err != nil {
		log.Warn("cannot rename %s to %s, trying to copy instead", src, dest)
	} else {
		return nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file %s, %v", src, err)
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("cannot create destination file %s, %v", dest, err)
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return fmt.Errorf("error while copying file, %v", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("cannot flush file contents, %v", err)
	}

	return nil
}

// RemoveOlderBackups keeps the most recent backups of a directory and deletes the old ones
func (f *Filesystem) RemoveOlderBackups(keep int) error {
	files, err := ioutil.ReadDir(f.SaveDir)
	if err != nil {
		return fmt.Errorf("cannot list contents of directory %s, %v", f.SaveDir, err)
	}

	count := len(files) - keep
	deleted := 0

	if count > 0 {
		for _, file := range files[:count] {
			err = os.Remove(file.Name())
			if err != nil {
				log.Error(0, "failed to remove file %s", file.Name())
			} else {
				deleted++
			}
		}

		log.Trace("deleted %d objects from %s", deleted, f.SaveDir)
	}

	return nil
}

// FindLatestBackup returns the most recent backup of the specified directory
func (f *Filesystem) FindLatestBackup() (string, error) {
	files, err := ioutil.ReadDir(f.SaveDir)
	if err != nil {
		return "", fmt.Errorf("cannot list contents of directory %s, %v", f.SaveDir, err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("cannot find a recent backup on %s", f.SaveDir)
	}

	return files[len(files)-1].Name(), nil
}

// Retrieve returns the path of the requested file
func (f *Filesystem) Retrieve(path string) (string, error) {
	return path, nil
}
