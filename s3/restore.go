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
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/urfave/cli"
)

type RestoreFunc func(path string) error

var DoRestore RestoreFunc
var saveDir = "/tmp"

func downloadFile(c *cli.Context, sess *session.Session, s3path string) (string, error) {
	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	file := path.Join(saveDir, path.Base(s3path))
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", file, err)
	}
	defer f.Close()

	// Upload the file to S3.
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(c.String("bucket")),
		Key:    aws.String(s3path),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download S3 object, %v", err)
	}

	fmt.Printf("file downloaded to %s\n", file)

	return file, nil
}

func findLatestBackup(sess *session.Session, c *cli.Context) (string, error) {
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
		return "", fmt.Errorf("couldn't list S3 objects, %v", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("cannot find a resent backup on s3://%s/%s",
			c.String("bucket"), c.String("prefix"))
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files[0], nil
}

func runRestoreTask(c *cli.Context) error {
	if DoRestore == nil {
		return fmt.Errorf("no restore function defined")
	}

	sess := getS3Session(c)
	s3path, err := findLatestBackup(sess, c)
	if err != nil {
		return fmt.Errorf("cannot find the latest backup, %v", err)
	}

	file, err := downloadFile(c, sess, s3path)
	if err != nil {
		return fmt.Errorf("cannot download S3 object %s, %v", s3path, err)
	}

	defer func() {
		err = os.Remove(file)
		if err != nil {
			log.Printf("cannot remove file %s, %v", file, err)
		}
	}()

	err = DoRestore(file)
	if err != nil {
		return fmt.Errorf("couldn't complete the restore, %v", err)
	}

	return nil
}

func RunRestore(c *cli.Context) error {
	runTask = runRestoreTask
	return runScheduler(c)
}
