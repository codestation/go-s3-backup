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
)

// ConsulConfig has the config options for the ConsulConfig service
type ConsulConfig struct {
	SaveDir string
}

// ConsulAppPath points to the consul binary location
var ConsulAppPath = "/bin/consul"

// Backup generates a tarball of the consul database and returns the path where is stored
func (c *ConsulConfig) Backup() (*BackupResults, error) {
	namePrefix := "consul-backup"
	filepath := generateFilename(c.SaveDir, namePrefix) + ".snap"
	args := []string{"snapshot", "save", filepath}

	app := CmdConfig{}

	if err := app.CmdRun(ConsulAppPath, args...); err != nil {
		return nil, fmt.Errorf("couldn't execute %s, %v", ConsulAppPath, err)
	}

	result := &BackupResults{[]BackupResult{{NamePrefix: namePrefix, Path: filepath}}}

	return result, nil
}

// Restore takes a GiteaConfig backup and restores it to the service
func (c *ConsulConfig) Restore(filepath string) error {
	args := []string{"snapshot", "restore", filepath}

	app := CmdConfig{}

	if err := app.CmdRun(ConsulAppPath, args...); err != nil {
		return fmt.Errorf("couldn't execute consul restore, %v", err)
	}

	return nil
}
