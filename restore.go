package main

import (
	"fmt"
	"os"
	"os/exec"
)

func gogsRestore(path string) error {
	cmd := exec.Command("gosu", "git", appPath, "restore", "--from", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "USER=git")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("couldn't execute %s, %v", appPath, err)
	}

	return nil
}
