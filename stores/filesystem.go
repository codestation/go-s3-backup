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
	"os"
	"path"

	log "gopkg.in/clog.v1"
)

type Filesystem struct {
	Path   string
	Rename bool
}

func (f *Filesystem) Store(filepath string, key string) error {
	dest := path.Join(f.Path, key)

	err := os.Rename(dest, filepath)
	if err != nil {
		log.Warn("cannot rename %s to %s, trying to copy instead", filepath, dest)
	} else {
		return nil
	}

	srcFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("cannot open file %s, %v", filepath, err)
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("cannot create file %s, %v", dest, err)
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

func (f *Filesystem) RemoveOlderBackups(prefix string, keep int) error {
	return fmt.Errorf("not implemented")
}

func (f *Filesystem) FindLatestBackup(prefix string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (f *Filesystem) Retrieve(path string) (string, error) {
	return path, nil
}
