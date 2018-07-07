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

type BackupFunc func(c *cli.Context) (string, error)

var DoBackup BackupFunc

func uploadFile(c *cli.Context, sess *session.Session, filepath string, key string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filepath, err)
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

func runBackupTask(c *cli.Context) error {
	if DoBackup == nil {
		return fmt.Errorf("no backup function defined")
	}

	filepath, err := DoBackup(c)
	if err != nil {
		return fmt.Errorf("couldn't complete the backup, %v", err)
	}

	defer func() {
		err = os.Remove(filepath)
		if err != nil {
			log.Printf("cannot remove file %s, %v", filepath, err)
		}
	}()

	sess := getS3Session(c)

	key := fmt.Sprintf("%s/%s", c.String("prefix"), path.Base(filepath))
	err = uploadFile(c, sess, filepath, key)
	if err != nil {
		return fmt.Errorf("couldn't upload the file to S3, %v", err)
	}

	if c.Int("max-backups") > 0 {
		if err = removeOlderBackups(sess, c); err != nil {
			return fmt.Errorf("couldn't remove olderbackups from S3, %v", err)
		}
	}

	return nil
}

func RunBackup(c *cli.Context) error {
	runTask = runBackupTask
	return runScheduler(c)
}
