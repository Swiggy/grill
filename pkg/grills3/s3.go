package grills3

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3 struct {
	minio *canned.Minio
}

func (gs *S3) Start(ctx context.Context) error {
	minio, err := canned.NewMinio(ctx)
	if err != nil {
		return err
	}
	gs.minio = minio
	return nil
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

func (gs *S3) Stop(ctx context.Context) error {
	return gs.minio.Container.Terminate(ctx)
}
