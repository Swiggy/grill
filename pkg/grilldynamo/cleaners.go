package grilldynamo

import (
	"bitbucket.org/swigy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (grilldynamo *GrillDynamo) DeleteTable(tableName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := grilldynamo.dynamo.Client.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		return err
	})
}
