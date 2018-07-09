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

package cmd

import (
	"fmt"
	"log"
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
)

type Task func(c *cli.Context) error

func BackupTask(c *cli.Context, service services.Service, store stores.Storer) error {
	filepath, err := service.Backup()
	if err != nil {
		return fmt.Errorf("couldn't make a gogs backup, %v", err)
	}

	defer func() {
		os.Remove(filepath)
	}()

	key := fmt.Sprintf("%s/%s", c.String("prefix"), path.Base(filepath))

	if err = store.Store(filepath, key); err != nil {
		return fmt.Errorf("couldn't upload file to S3, %v", err)
	}

	err = store.RemoveOlderBackups(c.String(c.String("prefix")), c.Int("max-backups"))
	if err != nil {
		fmt.Errorf("couldn't remove old backups from S3, %v", err)
	}

	return nil
}

func RestoreTask(c *cli.Context, service services.Service, store stores.Storer) error {
	var err error
	var s3path = c.String("s3path")

	if key := c.String("s3key"); key != "" {
		// restore directly from this S3 object
		s3path = key
	} else {
		// find the latest S3 object
		s3path, err = store.FindLatestBackup(c.String("prefix"))
		if err != nil {
			return fmt.Errorf("cannot find the latest backup, %v", err)
		}
	}

	filepath, err := store.Retrieve(s3path)
	if err != nil {
		return fmt.Errorf("cannot download S3 object %s, %v", s3path, err)
	}

	defer func() {
		if err := os.Remove(filepath); err != nil {
			log.Printf("cannot remove file %s, %v", filepath, err)
		}
	}()

	if err = service.Restore(filepath); err != nil {
		return fmt.Errorf("couldn't complete the restore, %v", err)
	}

	return nil
}

func RunScheduler(c *cli.Context, task Task) error {
	cr := cron.New()
	schedule := c.String("schedule")

	if schedule == "" || schedule == "none" {
		log.Printf("running job directly")
		return task(c)
	}

	log.Printf("starting scheduled backup jobs")
	timeoutchan := make(chan bool, 1)

	cr.AddFunc(schedule, func() {
		seconds := rand.Intn(c.Int("random-delay"))

		// run immediately is no delay is configured
		if seconds == 0 {
			if err := task(c); err != nil {
				log.Printf("failed to run scheduled task, %v", err)
			}
			return
		}

		log.Printf("waiting for %d seconds before starting scheduled job", seconds)

		select {
		case <-timeoutchan:
			log.Printf("random timeout cancelled")
			break
		case <-time.After(time.Duration(seconds) * time.Second):
			log.Printf("running scheduled task")

			if err := task(c); err != nil {
				log.Printf("failed to run scheduled task, %v", err)
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

	log.Printf("stopping scheduled task")
	cr.Stop()

	return nil
}
