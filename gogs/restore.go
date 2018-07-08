package gogs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func Restore(_ *cli.Context, filepath string) error {
	cmd := exec.Command("gosu", "git", AppPath, "restore", "--from", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	cmd.Env = append(env, "USER=git")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't execute %s, %v", AppPath, err)
	}

	return nil
}
