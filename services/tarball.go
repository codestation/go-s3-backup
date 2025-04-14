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
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archives"
)

// TarballConfig has the config options for the TarballConfig service
type TarballConfig struct {
	Name         string
	Path         string
	Compress     bool
	SaveDir      string
	BackupPerDir bool
	BackupDirs   []string
	ExcludeDirs  []string
}

// Backup creates a tarball of the specified directory
func (f *TarballConfig) Backup() (*BackupResults, error) {
	namePrefix := f.getNamePrefix("")
	if !f.BackupPerDir {
		filepath, err := f.backupFile("", namePrefix)
		if err != nil {
			return nil, err
		}

		return &BackupResults{Entries: []BackupResult{{
			NamePrefix: namePrefix,
			Path:       filepath,
		}}}, nil
	}

	files, err := os.ReadDir(path.Join(f.Path))
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

		filepath, err := f.backupFile(file.Name(), namePrefix)
		if err != nil {
			return nil, err
		}

		resultList = append(resultList, BackupResult{
			DirPrefix:  file.Name(),
			NamePrefix: namePrefix,
			Path:       filepath,
		})
	}

	result := &BackupResults{resultList}

	return result, nil
}

func (f *TarballConfig) getNamePrefix(basedir string) string {
	var name string
	switch {
	case f.Name != "":
		name = f.Name + "-backup"
	case basedir != "":
		name = path.Base(f.Path) + "_" + basedir + "-backup"
	default:
		name = path.Base(f.Path) + "-backup"
	}

	return name
}

func (f *TarballConfig) backupFile(basedir, namePrefix string) (string, error) {
	destPath := path.Join(f.SaveDir, basedir)
	filePath := generateFilename(destPath, namePrefix) + ".tar"

	format := archives.CompressedArchive{Archival: archives.Tar{}}

	if f.Compress {
		format.Compression = archives.Gz{}
		filePath += ".gz"
	}

	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return "", err
	}

	srcPath := path.Join(f.Path, basedir)

	ctx := context.TODO()

	options := &archives.FromDiskOptions{
		FollowSymlinks: false,
	}

	basePath := path.Base(srcPath)

	files, err := archives.FilesFromDisk(ctx, options, map[string]string{
		srcPath + string(os.PathSeparator): basePath,
	})
	if err != nil {
		return "", fmt.Errorf("cannot prepare tarball files on %s, %v", filePath, err)
	}

	cleanFilePath := filepath.Clean(filePath)

	out, err := os.Create(cleanFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = format.Archive(ctx, out, files)
	if err != nil {
		return "", fmt.Errorf("cannot create tarball on %s, %v", filePath, err)
	}

	return cleanFilePath, nil
}

// Restore extracts a tarball to the specified directory
func (f *TarballConfig) Restore(filepath string) error {
	err := removeDirectoryContents(f.Path)
	if err != nil {
		return fmt.Errorf("failed to empty directory contents before restoring: %v", err)
	}

	err = Unarchive(filepath, path.Dir(f.Path))
	if err != nil {
		return fmt.Errorf("cannot unpack backup: %v", err)
	}

	return nil
}

func Unarchive(source, destination string) error {
	ctx := context.TODO()

	// Open the source archive file
	archive, err := os.Open(source)
	if err != nil {
		return err
	}

	// Identify the archive file's format
	format, archiveReader, _ := archives.Identify(ctx, "", archive)

	dirMap := make(map[string]bool)

	// Check if the format is an extractor. If not, skip the archive file.
	extractor, ok := format.(archives.Extractor)

	if !ok {
		return nil
	}

	return extractor.Extract(ctx, archiveReader, func(_ context.Context, archiveFile archives.FileInfo) error {
		fileName := archiveFile.NameInArchive
		newPath := filepath.Join(destination, fileName)

		if archiveFile.IsDir() {
			dirMap[newPath] = true

			return os.MkdirAll(newPath, 0o755) // #nosec
		}

		fileDir := filepath.Dir(newPath)
		_, seenDir := dirMap[fileDir]

		if !seenDir {
			dirMap[fileDir] = true

			_ = os.MkdirAll(fileDir, 0o755) // #nosec
		}

		cleanNewPath := filepath.Clean(newPath)

		newFile, err := os.OpenFile(cleanNewPath,
			os.O_CREATE|os.O_WRONLY,
			archiveFile.Mode())
		if err != nil {
			return err
		}
		defer newFile.Close()

		archiveFileTemp, err := archiveFile.Open()
		if err != nil {
			return err
		}
		defer archiveFileTemp.Close()

		_, err = io.Copy(newFile, archiveFileTemp)

		return err
	})
}
