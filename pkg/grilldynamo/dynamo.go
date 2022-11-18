package grilldynamo

import (
	"context"
	"github.com/Swiggy/grill/canned"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Dynamo struct {
	dynamo *canned.DynamoDB
}

func (gd *Dynamo) Start(ctx context.Context) error {
	dynamo, err := canned.NewDynamoDB(ctx)
	if err != nil {
		return err
	}
	gd.dynamo = dynamo
	return nil
}

func (gd *Dynamo) Client() dynamodbiface.DynamoDBAPI {
	return gd.dynamo.Client
}

func (gd *Dynamo) Host() string {
	return gd.dynamo.Host
}

func (gd *Dynamo) Port() string {
	return gd.dynamo.Port
}

func (gd *Dynamo) Region() string {
	return gd.dynamo.Region
}

func (gd *Dynamo) AccessKey() string {
	return gd.dynamo.AccessKey
}

func (gd *Dynamo) SecretKey() string {
	return gd.dynamo.SecretKey
}

func (gd *Dynamo) Stop(ctx context.Context) error {
	return gd.dynamo.Container.Terminate(ctx)
}
