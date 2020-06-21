package grilldynamo

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type GrillDynamo struct {
	dynamo *canned.DynamoDB
}

func Start() (*GrillDynamo, error) {
	dynamo, err := canned.NewDynamoDB(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GrillDynamo{
		dynamo: dynamo,
	}, nil
}

func (grilldynamo *GrillDynamo) Client() dynamodbiface.DynamoDBAPI {
	return grilldynamo.dynamo.Client
}

func (grilldynamo *GrillDynamo) Host() string {
	return grilldynamo.dynamo.Host
}

func (grilldynamo *GrillDynamo) Port() string {
	return grilldynamo.dynamo.Port
}

func (grilldynamo *GrillDynamo) Region() string {
	return grilldynamo.dynamo.Region
}

func (grilldynamo *GrillDynamo) AccessKey() string {
	return grilldynamo.dynamo.AccessKey
}

func (grilldynamo *GrillDynamo) SecretKey() string {
	return grilldynamo.dynamo.SecretKey
}

func (grilldynamo *GrillDynamo) Stop() error {
	return grilldynamo.dynamo.Container.Terminate(context.TODO())
}
