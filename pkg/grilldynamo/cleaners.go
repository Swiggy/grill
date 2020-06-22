package grilldynamo

import (
	"bitbucket.org/swigy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (gd *Dynamo) DeleteTable(tableName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gd.dynamo.Client.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		return err
	})
}
