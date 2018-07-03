package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
	"math/rand"
	"path"
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

func uploadFile(c *cli.Context, filename string, key string) error {
	s3Config := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", "default"),
		Endpoint:    aws.String(c.String("endpoint")),
		Region:      aws.String(c.String("region")),
	}

	if c.Bool("force-path-style") {
		s3Config.S3ForcePathStyle = aws.Bool(true)
	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(s3Config))

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

	key := fmt.Sprintf("%s/%s", c.String("prefix"), filename)
	err = uploadFile(c, path.Join(saveDir, filename), key)
	if err != nil {
		return fmt.Errorf("couldn't upload the file to S3, %v", err)
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
