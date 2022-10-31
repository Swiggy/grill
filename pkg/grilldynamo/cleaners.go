package grilldynamo

import (
	"github.com/Swiggy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (gd *Dynamo) DeleteItem(tableName string, key map[string]*dynamodb.AttributeValue) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gd.dynamo.Client.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key:       key,
		})
		return err
	})
}

func (gd *Dynamo) DeleteTable(tableName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gd.dynamo.Client.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		return err
	})
}
