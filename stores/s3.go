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
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "gopkg.in/clog.v1"
)

// S3 has the config options for the S3 service
type S3 struct {
	Endpoint          string
	Region            string
	Bucket            string
	AccessKey         string
	ClientSecret      string
	Prefix            string
	ForcePathStyle    bool
	RemoveAfterUpload bool
	SaveDir           string
	retrievedFile     string
}

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

// Store saves a file to a remote S3 service
func (s *S3) Store(filepath string, filename string) error {
	uploader := s3manager.NewUploader(s.newSession())

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filepath, err)
	}

	if s.RemoveAfterUpload {
		defer func() {
			log.Info("removing source file %s", filepath)
			if err = os.Remove(filepath); err != nil {
				log.Warn("cannot remove file %s, %v", filepath, err)
			}
		}()
	}

	key := path.Clean(path.Join(s.Prefix, filename))

	// Upload the file to S3.
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	log.Trace("file uploaded to %s\n", res.Location)

	return nil
}

// RemoveOlderBackups keeps the most recent backups of the S3 service and deletes the old ones
func (s *S3) RemoveOlderBackups(keep int) error {
	svc := s3.New(s.newSession())

	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		// make sure that the prefix ends with "/"
		Prefix: aws.String(path.Clean(s.Prefix) + "/"),
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
	count := len(files) - keep

	if count > 0 {
		var items s3.Delete
		var objs = make([]*s3.ObjectIdentifier, count)

		for i, file := range files[:count] {
			objs[i] = &s3.ObjectIdentifier{Key: aws.String(file)}
			log.Trace("marked to delete: s3://%s/%s", s.Bucket, file)
		}

		items.SetObjects(objs)

		out, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(s.Bucket),
			Delete: &items})

		if err != nil {
			return fmt.Errorf("couldn't delete the S3 objects, %v", err)
		}

		log.Trace("deleted %d objects from S3", len(out.Deleted))
	}

	return nil
}

// FindLatestBackup returns the most recent backup of the S3 store
func (s *S3) FindLatestBackup() (string, error) {
	svc := s3.New(s.newSession())

	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		// make sure that the prefix ends with "/"
		Prefix: aws.String(path.Clean(s.Prefix) + "/"),
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
			s.Bucket, s.Prefix)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files[0], nil
}

// Retrieve downloads a S3 object to the local filesystem
func (s *S3) Retrieve(s3path string) (string, error) {
	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(s.newSession())

	filepath := path.Join(s.SaveDir, path.Base(s3path))
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

	log.Trace("file downloaded to %s\n", filepath)
	s.retrievedFile = filepath

	return filepath, nil
}

// Close deinitializes the store (remove downloaded file)
func (s *S3) Close() {
	if s.retrievedFile != "" {
		if err := os.Remove(s.retrievedFile); err != nil {
			log.Warn("cannot remove file %s", s.retrievedFile)
		}

		s.retrievedFile = ""
	}
}
