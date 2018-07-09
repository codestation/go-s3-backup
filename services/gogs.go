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
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"
)

type Gogs struct {
	ConfigPath string
	DataPath   string
}

var GogsAppPath = "/app/gogs/gogs"

func getEnvInt(key string, def int) int {
	value := os.Getenv(key)

	if value != "" {
		val, err := strconv.Atoi(value)
		if err == nil {
			return val
		} else {
			log.Printf("cannot parse env key %s with value %s", key, value)
		}
	}

	return def
}

func (g *Gogs) Backup() (string, error) {
	filepath := fmt.Sprintf("gogs-backup-%s.zip", time.Now().Format("20060102150405"))

	args := []string{
		"backup",
		"--target", SaveDir,
		"--archive-name", filepath,
	}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	cmd := exec.Command(GogsAppPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// only switch user when running as root
	if os.Geteuid() == 0 {
		// run process as git user
		uid := uint32(getEnvInt("PUID", 1000))
		gid := uint32(getEnvInt("PGID", 1000))

		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	} else {
		log.Printf("not runnign as root, starting %s with UID %d", GogsAppPath, os.Geteuid())
	}

	env := os.Environ()
	home := fmt.Sprintf("HOME=%s", path.Join(g.DataPath, "git"))
	cmd.Env = append(env, "USER=git", home)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", GogsAppPath, err)
	}

	return path.Join(SaveDir, filepath), nil
}

func (g *Gogs) Restore(filepath string) error {
	args := []string{
		"restore",
		"--from", filepath,
	}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	cmd := exec.Command(GogsAppPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// only switch user when running as root
	if os.Geteuid() == 0 {
		// run process as git user
		uid := uint32(getEnvInt("PUID", 1000))
		gid := uint32(getEnvInt("PGID", 1000))

		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	} else {
		log.Printf("not runnign as root, starting %s with UID %d", GogsAppPath, os.Geteuid())
	}

	env := os.Environ()
	home := fmt.Sprintf("HOME=%s", path.Join(g.DataPath, "git"))
	cmd.Env = append(env, "USER=git", home)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't execute gogs restore, %v", err)
	}

	return nil
}
