package grills3

import (
	"bitbucket.org/swigy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (gs *S3) DeleteBucket(bucketName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gs.Client().DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		return err
	})
}

func (gs *S3) DeleteAllFiles(bucketName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		output, err := gs.Client().ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		for _, object := range output.Contents {
			_, err := gs.Client().DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    object.Key,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
