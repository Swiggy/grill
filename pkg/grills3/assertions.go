package grills3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gs *S3) AssertFileExists(bucket, key string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		_, err := gs.Client().GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("no object filed for key=%v, in bucket=%v, error=%v", key, bucket, err)
		}
		return nil
	})
}
