package grills3

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3 struct {
	minio *canned.Minio
}

func Start() (*S3, error) {
	minio, err := canned.NewMinio(context.TODO())
	if err != nil {
		return nil, err
	}

	return &S3{
		minio: minio,
	}, nil
}

func (gs *S3) Host() string {
	return gs.minio.Host
}

func (gs *S3) Port() string {
	return gs.minio.Port
}

func (gs *S3) Client() s3iface.S3API {
	return gs.minio.Client
}

func (gs *S3) Region() string {
	return gs.minio.Region
}

func (gs *S3) AccessKey() string {
	return gs.minio.AccessKey
}

func (gs *S3) SecretKey() string {
	return gs.minio.SecretKey
}

func (gs *S3) Stop() error {
	return gs.minio.Container.Terminate(context.Background())
}
