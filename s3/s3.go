package s3

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/externalservices"
)

type S3 struct {
	svc    *awss3.S3
	bucket *string
	path   *string
}

func NewS3(s3config externalservices.S3Config) S3 {
	svc := awss3.New(session.New(), &aws.Config{
		Region:      aws.String(s3config.Region),
		Endpoint:    aws.String(s3config.Endpoint),
		Credentials: credentials.NewStaticCredentials(s3config.AccessKeyID, s3config.SecretAccessKey, ""),
	})
	return S3{
		svc:    svc,
		bucket: aws.String(s3config.Bucket),
		path:   aws.String(s3config.Path),
	}
}

func (s *S3) FetchNewestFile() ([]byte, error) {
	lastModifiedOfNewestFile := *new(time.Time)
	filenameOfNewestFile := ""

	err := s.svc.ListObjectsV2Pages(&awss3.ListObjectsV2Input{
		Prefix: s.path,
		Bucket: s.bucket,
	}, func(page *awss3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.Contents {
			if object.LastModified.After(lastModifiedOfNewestFile) {
				filenameOfNewestFile = *object.Key
				lastModifiedOfNewestFile = *object.LastModified
			}
		}
		return true
	})
	if err != nil {
		return []byte{}, err
	}

	if filenameOfNewestFile == "" {
		return []byte{}, fmt.Errorf("No previous dataset found")
	}

	out, err := s.svc.GetObject(&awss3.GetObjectInput{
		Bucket: s.bucket,
		Key:    &filenameOfNewestFile,
	})
	if err != nil {
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(out.Body)
	return buf.Bytes(), err
}

func (s *S3) PutDataset(filename string, dataset []byte) error {
	r := bytes.NewReader(dataset)
	_, err := s.svc.PutObject(&awss3.PutObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(fmt.Sprintf("%s/%s", *s.path, filename)),
		Body:   r,
	})

	if err != nil {
		return err
	}

	return nil
}
