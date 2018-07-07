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

package s3

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

type taskFunc func(c *cli.Context) error

var runTask taskFunc

func getS3Session(c *cli.Context) *session.Session {
	s3Config := &aws.Config{
		Credentials:      credentials.NewSharedCredentials("", "default"),
		Endpoint:         aws.String(c.String("endpoint")),
		Region:           aws.String(c.String("region")),
		S3ForcePathStyle: aws.Bool(c.Bool("force-path-style")),
	}

	return session.Must(session.NewSession(s3Config))
}

func runScheduler(c *cli.Context) error {
	cr := cron.New()
	schedule := c.String("schedule")

	if schedule == "" || schedule == "none" {
		log.Printf("running scheduler job")

		if runTask == nil {
			panic("runTask function isn't defined")
		}

		return runTask(c)
	}

	log.Printf("starting scheduled backup jobs")
	timeoutchan := make(chan bool, 1)

	cr.AddFunc(schedule, func() {
		seconds := rand.Intn(c.Int("random-delay"))

		// run immediately is no delay is configured
		if seconds == 0 {
			if err := runTask(c); err != nil {
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
			log.Printf("running backup job")

			if runTask == nil {
				panic("runTask function isn't defined")
			}

			if err := runTask(c); err != nil {
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

	log.Printf("stopping scheduled job")
	cr.Stop()

	return nil
}
