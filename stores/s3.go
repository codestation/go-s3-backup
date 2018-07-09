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

package stores

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Endpoint       string
	Region         string
	Bucket         string
	AccessKey      string
	ClientSecret   string
	ForcePathStyle bool
}

var SaveDir = "/tmp"

func (s *S3) newSession() *session.Session {
	var creds *credentials.Credentials

	if s.AccessKey != "" && s.ClientSecret != "" {
		creds = credentials.NewStaticCredentials(s.AccessKey, s.ClientSecret, "")
	} else {
		creds = credentials.NewSharedCredentials("", "default")
	}

	s3Config := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(s.Endpoint),
		Region:           aws.String(s.Region),
		S3ForcePathStyle: aws.Bool(s.ForcePathStyle),
	}

	return session.Must(session.NewSession(s3Config))
}

func (s *S3) Store(filepath string, key string) error {
	uploader := s3manager.NewUploader(s.newSession())

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filepath, err)
	}

	defer func() {
		os.Remove(filepath)
	}()

	// Upload the file to S3.
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to %s\n", res.Location)

	return nil
}

func (s *S3) RemoveOlderBackups(prefix string, keep int) error {
	svc := s3.New(s.newSession())

	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		// make sure that the prefix ends with "/"
		Prefix: aws.String(path.Clean(prefix) + "/"),
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
	count := len(files) - keep
	var objs = make([]*s3.ObjectIdentifier, count)

	for i, file := range files[:count] {
		objs[i] = &s3.ObjectIdentifier{Key: aws.String(file)}
		log.Printf("marked to delete: %s", file)
	}

	items.SetObjects(objs)

	out, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s.Bucket),
		Delete: &items})

	if err != nil {
		return fmt.Errorf("couldn't delete the S3 objects, %v", err)
	} else {
		fmt.Printf("deleted %d objects from S3", len(out.Deleted))
	}

	return nil
}

func (s *S3) FindLatestBackup(prefix string) (string, error) {
	svc := s3.New(s.newSession())

	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		// make sure that the prefix ends with "/"
		Prefix: aws.String(path.Clean(prefix) + "/"),
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
		return "", fmt.Errorf("cannot find a recent backup on s3://%s/%s",
			s.Bucket, prefix)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files[0], nil
}

func (s *S3) Retrieve(s3path string) (string, error) {
	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(s.newSession())

	filepath := path.Join(SaveDir, path.Base(s3path))
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filepath, err)
	}

	defer f.Close()

	// Upload the file to S3.
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3path),
	})

	if err != nil {
		return "", fmt.Errorf("failed to download S3 object, %v", err)
	}

	fmt.Printf("file downloaded to %s\n", filepath)

	return filepath, nil
}
