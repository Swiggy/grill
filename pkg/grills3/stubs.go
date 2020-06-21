package grills3

import (
	"os"

	"bitbucket.org/swigy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (gs *GrillS3) CreateBucket(bucketName string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := gs.Client().CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		return err
	})
}

func (gs *GrillS3) UploadFile(bucket, key, filepath string) grill.Stub {
	return grill.StubFunc(func() error {
		file, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = gs.Client().PutObject(&s3.PutObjectInput{
			Body:   file,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		return err
	})
}
