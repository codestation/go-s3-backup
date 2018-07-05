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
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

var saveDir = "/tmp"

func gogsBackup() (string, error) {
	filename := fmt.Sprintf("gogs-backup-%s.zip", time.Now().Format("20060102150405"))
	cmd := exec.Command("gosu", "git", appPath, "backup", "--target", saveDir, "--archive-name", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "USER=git")
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", appPath, err)
	}

	return filename, nil
}

func getS3Session(c *cli.Context) *session.Session {
	s3Config := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", "default"),
		Endpoint:    aws.String(c.String("endpoint")),
		Region:      aws.String(c.String("region")),
	}

	if c.Bool("force-path-style") {
		s3Config.S3ForcePathStyle = aws.Bool(true)
	}

	// The session the S3 Uploader will use
	return session.Must(session.NewSession(s3Config))
}

func uploadFile(c *cli.Context, sess *session.Session, filename string, key string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.String("bucket")),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to %s\n", result.Location)

	return nil
}

func removeOlderBackups(sess *session.Session, c *cli.Context) error {
	svc := s3.New(sess)

	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(c.String("bucket")),
		// make sure that the prefix ends with "/"
		Prefix: aws.String(path.Clean(c.String("prefix")) + "/"),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {

		for _, obj := range p.Contents {
			if !strings.HasSuffix(*obj.Key, "/") {
				files = append(files, aws.StringValue(obj.Key))
			}
		}
		return true
	})

	if err != nil {
		return fmt.Errorf("couldn't list S3 objects, %v", err)
	}

	sort.Strings(files)

	var items s3.Delete
	count := len(files) - c.Int("max-backups")
	var objs = make([]*s3.ObjectIdentifier, count)

	for i, file := range files[:count] {
		objs[i] = &s3.ObjectIdentifier{Key: aws.String(file)}
		log.Printf("marked to delete: %s", file)
	}

	items.SetObjects(objs)

	out, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(c.String("bucket")),
		Delete: &items})

	if err != nil {
		return fmt.Errorf("couldn't delete the S3 objects, %v", err)
	} else {
		fmt.Printf("deleted %d objects from S3", len(out.Deleted))
	}

	return nil
}

func runTask(c *cli.Context) error {
	filename, err := gogsBackup()
	if err != nil {
		return fmt.Errorf("couldn't complete the backup, %v", err)
	}

	defer func() {
		err = os.Remove(path.Join(saveDir, filename))
		if err != nil {
			log.Printf("cannot remove file %s, %v", filename, err)
		}
	}()

	sess := getS3Session(c)

	key := fmt.Sprintf("%s/%s", c.String("prefix"), filename)
	err = uploadFile(c, sess, path.Join(saveDir, filename), key)
	if err != nil {
		return fmt.Errorf("couldn't upload the file to S3, %v", err)
	}

	if c.Int("max-backups") > 0 {
		err = removeOlderBackups(sess, c)
		if err != nil {
			return fmt.Errorf("couldn't remove olderbackups from S3, %v", err)
		}
	}

	return nil
}

func run(c *cli.Context) error {
	cr := cron.New()
	schedule := c.String("schedule")

	if schedule == "" || schedule == "none" {
		log.Printf("running backup job")

		return runTask(c)
	}

	log.Printf("starting scheduled backup jobs")
	timeoutchan := make(chan bool, 1)

	cr.AddFunc(schedule, func() {
		minutes := rand.Intn(60)
		log.Printf("waiting for %d minutes before starting backup job", minutes)

		select {
		case <-timeoutchan:
			log.Printf("random timeout cancelled")
			break
		case <-time.After(time.Duration(minutes) * time.Minute):
			log.Printf("running backup job")
			runTask(c)
			break
		}
	})
	cr.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	timeoutchan <- true
	close(timeoutchan)

	log.Printf("stopping scheduled jobs")
	cr.Stop()

	return nil
}
