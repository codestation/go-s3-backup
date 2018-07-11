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
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"megpoid.xyz/go/go-s3-backup/services"
	"megpoid.xyz/go/go-s3-backup/stores"

	"github.com/robfig/cron"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
)

type Task func(c *cli.Context) error

func getService(c *cli.Context, service string) services.Service {
	var serv services.Service
	switch service {
	case "gogs":
		serv = NewGogsConfig(c)
	case "mysql":
		serv = NewMysqlConfig(c)
	case "postgres":
		serv = NewPostgresConfig(c)
	case "tarball":
		serv = NewTarballConfig(c)
	default:
		log.Fatal(0, "unsupported service: %s", service)
	}

	return serv
}

func getStore(c *cli.Context, store string) stores.Storer {
	var serv stores.Storer
	switch store {
	case "s3":
		serv = NewS3Config(c)
	case "filesystem":
		serv = NewFilesystemConfig(c)
	default:
		log.Fatal(0, "unsupported store: %s", store)
	}

	return serv
}

func runTask(c *cli.Context, command string, serviceName string, storeName string) error {
	service := getService(c, serviceName)
	store := getStore(c, storeName)

	switch command {
	case "backup":
		return RunScheduler(c, func(c *cli.Context) error {
			return BackupTask(c, service, store)
		})
	case "restore":
		return RunScheduler(c, func(c *cli.Context) error {
			return RestoreTask(c, service, store)
		})
	default:
		log.Fatal(0, "unsupported command: %s", command)
	}
	return nil
}

func BackupTask(c *cli.Context, service services.Service, store stores.Storer) error {
	filepath, err := service.Backup()
	if err != nil {
		return fmt.Errorf("service backup failed: %v", err)
	}

	log.Trace("backup saved to %s", filepath)

	filename := path.Base(filepath)

	if err = store.Store(filepath, filename); err != nil {
		return fmt.Errorf("couldn't upload file to S3, %v", err)
	}

	err = store.RemoveOlderBackups(c.GlobalInt("max-backups"))
	if err != nil {
		return fmt.Errorf("couldn't remove old backups from S3, %v", err)
	}

	return nil
}

func RestoreTask(c *cli.Context, service services.Service, store stores.Storer) error {
	var err error
	var s3key string

	if key := c.GlobalString("s3key"); key != "" {
		// restore directly from this S3 object
		s3key = key
	} else {
		// find the latest S3 object
		s3key, err = store.FindLatestBackup()
		if err != nil {
			return fmt.Errorf("cannot find the latest backup, %v", err)
		}
	}

	filepath, err := store.Retrieve(s3key)
	if err != nil {
		return fmt.Errorf("cannot download S3 object %s, %v", s3key, err)
	}

	log.Trace("backup retrieved to %s", filepath)

	defer func() {
		if err := os.Remove(filepath); err != nil {
			log.Warn("cannot remove file %s, %v", filepath, err)
		}
	}()

	if err = service.Restore(filepath); err != nil {
		return fmt.Errorf("service restore failed: %v", err)
	}

	return nil
}

func RunScheduler(c *cli.Context, task Task) error {
	cr := cron.New()
	schedule := c.GlobalString("schedule")

	if schedule == "" || schedule == "none" {
		log.Trace("running job directly")
		return task(c)
	}

	log.Trace("starting scheduled backup jobs")
	timeoutchan := make(chan bool, 1)

	cr.AddFunc(schedule, func() {
		seconds := rand.Intn(c.GlobalInt("random-delay"))

		// run immediately is no delay is configured
		if seconds == 0 {
			if err := task(c); err != nil {
				log.Error(0, "failed to run scheduled task, %v", err)
			}
			return
		}

		log.Trace("waiting for %d seconds before starting scheduled job", seconds)

		select {
		case <-timeoutchan:
			log.Trace("random timeout cancelled")
			break
		case <-time.After(time.Duration(seconds) * time.Second):
			log.Trace("running scheduled task")

			if err := task(c); err != nil {
				log.Error(0, "failed to run scheduled task, %v", err)
			}
			break
		}
	})
	cr.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	timeoutchan <- true
	close(timeoutchan)

	log.Trace("stopping scheduled task")
	cr.Stop()

	return nil
}
