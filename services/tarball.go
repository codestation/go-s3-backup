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

package services

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"path"
)

// TarballConfig has the config options for the TarballConfig service
type TarballConfig struct {
	Name         string
	Path         string
	Compress     bool
	SaveDir      string
	Prefix       string
	BackupPerDir bool
	BackupDirs   []string
	ExcludeDirs  []string
}

// Backup creates a tarball of the specified directory
func (f *TarballConfig) Backup() (*BackupResults, error) {
	if !f.BackupPerDir {
		filepath, err := f.backupFile("")
		if err != nil {
			return nil, err
		}

		return &BackupResults{Entries: []BackupResult{{
			Filenames: []string{filepath},
		}}}, nil
	}

	files, err := ioutil.ReadDir(path.Join(f.Path))
	if err != nil {
		return nil, err
	}

	var resultList []BackupResult

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		found := false
		for _, entry := range f.ExcludeDirs {
			if entry == file.Name() {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if len(f.BackupDirs) > 0 {
			found := false
			for _, entry := range f.BackupDirs {
				if entry == path.Base(file.Name()) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filepath, err := f.backupFile(file.Name())
		if err != nil {
			return nil, err
		}

		resultList = append(resultList, BackupResult{
			Prefix:    file.Name(),
			Filenames: []string{filepath},
		})
	}

	result := &BackupResults{resultList}

	return result, nil
}

func (f *TarballConfig) backupFile(basedir string) (string, error) {
	var name string
	if f.Name != "" {
		name = f.Name + "-backup"
	} else {
		name = path.Base(f.Path) + "_" + basedir + "-backup"
	}

	destPath := path.Join(f.SaveDir, basedir, f.Prefix)
	filepath := generateFilename(destPath, name) + ".tar"

	if f.Compress {
		filepath += ".gz"
	}

	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		return "", err
	}

	srcPath := path.Join(f.Path, basedir, f.Prefix)

	err = archiver.Archive([]string{srcPath}, filepath)
	if err != nil {
		return "", fmt.Errorf("cannot create tarball on %s, %v", filepath, err)
	}

	return filepath, nil
}

// Restore extracts a tarball to the specified directory
func (f *TarballConfig) Restore(filepath string) error {
	err := removeDirectoryContents(f.Path)
	if err != nil {
		return fmt.Errorf("failed to empty directory contents before restoring: %v", err)
	}

	err = archiver.Unarchive(filepath, path.Dir(f.Path))
	if err != nil {
		return fmt.Errorf("cannot unpack backup: %v", err)
	}

	return nil
}
