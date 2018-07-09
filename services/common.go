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
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var SaveDir = "/tmp"

type Service interface {
	Backup() (string, error)
	Restore(path string) error
}

func CompressAppOutput(cmd *exec.Cmd, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("cannot open file %s, %v", filepath, err)
	}

	defer f.Close()

	pr, pw := io.Pipe()
	gzW := gzip.NewWriter(pw)

	cmd.Stdout = gzW

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("couldn't execute %s, %v", cmd.Args[0], err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
		gzW.Close()
		pw.Close()
	}()

	_, err = io.Copy(f, pr)
	if err != nil {
		return fmt.Errorf("couldn't pipe command stdout to file, %v", err)
	}

	return nil
}

func ReadFileToInput(cmd *exec.Cmd, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("cannot open file %s, %v", filepath, err)
	}

	defer f.Close()

	pr, pw := io.Pipe()
	var gzR io.Reader

	if strings.HasSuffix(filepath, ".gz") {
		gzR, err = gzip.NewReader(pr)
		if err != nil {
			return fmt.Errorf("cannot create gzip reader, %v", err)
		}
	} else {
		gzR = pr
	}

	cmd.Stdin = gzR

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("couldn't execute %s, %v", cmd.Args[0], err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
		pw.Close()
	}()

	_, err = io.Copy(pw, f)
	if err != nil {
		return fmt.Errorf("couldn't pipe file contents to stdin, %v", err)
	}

	return nil
}
