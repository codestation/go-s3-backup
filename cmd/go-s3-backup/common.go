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

package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.megpoid.dev/go-s3-backup/services"
	"go.megpoid.dev/go-s3-backup/stores"
)

type task func(c *cli.Context) error

func getService(c *cli.Context, service string) services.Service {
	var config services.Service
	switch service {
	case "mysql":
		config = newMysqlConfig(c)
	case "postgres":
		config = newPostgresConfig(c)
	case "tarball":
		config = newTarballConfig(c)
	default:
		slog.Error("Unsupported service", "service", service)
		os.Exit(1)
	}

	return config
}

func getStore(c *cli.Context, store string) stores.Storer {
	var config stores.Storer
	switch store {
	case "s3":
		config = newS3Config(c)
	case "filesystem":
		config = newFilesystemConfig(c)
	default:
		slog.Error("Unsupported store", "store", store)
		os.Exit(1)
	}

	return config
}

func runTask(c *cli.Context, command string, serviceName string, storeName string) error {
	service := getService(c, serviceName)
	store := getStore(c, storeName)

	switch command {
	case "backup":
		return runScheduler(c, func(c *cli.Context) error {
			return backupTask(c, service, store)
		})
	case "restore":
		return runScheduler(c, func(c *cli.Context) error {
			return restoreTask(c, service, store)
		})
	default:
		slog.Error("Unsupported command", "command", command)
		os.Exit(1)
	}
	return nil
}

func backupTask(c *cli.Context, service services.Service, store stores.Storer) error {
	results, err := service.Backup()
	if err != nil {
		return fmt.Errorf("service backup failed: %v", err)
	}

	for _, result := range results.Entries {
		slog.Debug("Backup saved", "basedir", result.DirPrefix, "path", result.Path)
		filename := path.Base(result.Path)
		if err = store.Store(result.Path, result.DirPrefix, filename); err != nil {
			return fmt.Errorf("couldn't upload file to store: %v", err)
		}

		err = store.RemoveOlderBackups(result.DirPrefix, result.NamePrefix, c.Int("max-backups"))
		if err != nil {
			return fmt.Errorf("couldn't remove old backups from store: %v", err)
		}
	}

	return nil
}

func restoreTask(c *cli.Context, service services.Service, store stores.Storer) error {
	var err error
	var filename string

	if key := c.String("restore-file"); key != "" {
		// restore directly from this file
		filename = key
	} else {
		// find the latest file in the store
		filename, err = store.FindLatestBackup("", c.String("restore-prefix"))
		if err != nil {
			return fmt.Errorf("cannot find the latest backup: %v", err)
		}
	}

	filepath, err := store.Retrieve(filename)
	if err != nil {
		return fmt.Errorf("cannot download file %s: %v", filename, err)
	}

	defer store.Close()

	if err = service.Restore(filepath); err != nil {
		return fmt.Errorf("service restore failed: %v", err)
	}

	return nil
}

func runScheduler(c *cli.Context, task task) error {
	cr := cron.New()
	schedule := c.String("schedule")

	if schedule == "" || schedule == "none" {
		slog.Debug("Running task directly")
		return task(c)
	}

	slog.Debug("Starting scheduled backup task")
	timeoutchan := make(chan bool, 1)

	_, err := cr.AddFunc(schedule, func() {
		delay := c.Int("random-delay")
		if delay <= 0 {
			slog.Warn("Schedule random delay was set to a number <= 0, using 1 as default")
			delay = 1
		}

		seconds := rand.Intn(delay)

		// run immediately is no delay is configured
		if seconds == 0 {
			if err := task(c); err != nil {
				slog.Error("Failed to run scheduled task", "error", err)
			}
			return
		}

		slog.Debug("Waiting before starting scheduled job", "seconds", seconds)

		select {
		case <-timeoutchan:
			slog.Debug("Random timeout cancelled")
			break
		case <-time.After(time.Duration(seconds) * time.Second):
			slog.Debug("Running scheduled task")

			if err := task(c); err != nil {
				slog.Error("Failed to run scheduled task", "error", err)
			}
			break
		}
	})
	if err != nil {
		return err
	}

	cr.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	timeoutchan <- true
	close(timeoutchan)

	slog.Debug("Received signal, stopping scheduled task")
	ctx := cr.Stop()
	<-ctx.Done()

	return nil
}

func fileOrString(c *cli.Context, name string) string {
	if filepath := c.String(name + "-file"); filepath != "" {
		f, err := os.Open(filepath)
		if err != nil {
			slog.Error("Cannot open file", "filepath", filepath, "error", err)
			return ""
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				slog.Error("Cannot close file", "filepath", filepath, "error", err)
			}
		}(f)

		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			return scanner.Text()
		}

		slog.Warn("Empty file", "filepath", filepath)
		return ""
	}

	return c.String(name)
}

func applyConfigValues(flags []cli.Flag) cli.BeforeFunc {
	return func(c *cli.Context) error {
		config := c.App.Metadata["config"]
		if config != nil {
			cfg, ok := config.(altsrc.InputSourceContext)
			if ok {
				return altsrc.ApplyInputSourceValues(c, cfg, flags)
			}

			return fmt.Errorf("invalid config type for metadata")
		}

		return nil
	}
}
