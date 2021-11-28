package grilldynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/swiggy-private/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type testVal struct {
	Value string
}

type testItem struct {
	PartitionKey string
	Name         testVal
}

var (
	testTableRequest = &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
				AttributeName: aws.String("PartitionKey"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				KeyType:       aws.String(dynamodb.KeyTypeHash),
				AttributeName: aws.String("PartitionKey"),
			},
		},
		TableName: aws.String("testTable"),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}

	testDBItem, _ = dynamodbattribute.MarshalMap(&testItem{PartitionKey: "test", Name: testVal{Value: "test"}})
)

func Test_GrillDynamo(t *testing.T) {
	helper := &Dynamo{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting dynamo grill, error=%v", err)
		return
	}

	tableName := "testTable"

	tests := []grill.TestCase{
		{
			Name: "TestDynamo_SeedCount",
			Stubs: []grill.Stub{
				helper.CreateTable(testTableRequest),
				helper.SeedDataFromFile(tableName, "test_data/db.seed"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertScanCount(&dynamodb.ScanInput{TableName: aws.String(tableName)}, 3),
				helper.AssertItem(&dynamodb.GetItemInput{
					TableName: aws.String(tableName),
					Key:       map[string]*dynamodb.AttributeValue{"PartitionKey": {S: aws.String("3")}},
				}, testItem{PartitionKey: "3", Name: testVal{Value: "three"}}),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteTable(tableName),
			},
		},
		{
			Name: "TestDynamo_PutItem",
			Stubs: []grill.Stub{
				helper.CreateTable(testTableRequest),
				helper.PutItem(tableName, testDBItem),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertScanCount(&dynamodb.ScanInput{TableName: aws.String(tableName)}, 1),
				helper.AssertItem(&dynamodb.GetItemInput{
					TableName: aws.String(tableName),
					Key:       map[string]*dynamodb.AttributeValue{"PartitionKey": {S: aws.String("test")}},
				}, testItem{PartitionKey: "test", Name: testVal{Value: "test"}}),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteItem(tableName, map[string]*dynamodb.AttributeValue{"PartitionKey": {S: aws.String("test")}}),
				helper.DeleteTable(tableName),
			},
		},
		{
			Name: "TestDynamo_ItemNotFound",
			Stubs: []grill.Stub{
				helper.CreateTable(testTableRequest),
				helper.SeedDataFromFile(tableName, "test_data/db.seed"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertScanCount(&dynamodb.ScanInput{TableName: aws.String(tableName)}, 3),
				helper.AssertItem(&dynamodb.GetItemInput{
					TableName: aws.String(tableName),
					Key:       map[string]*dynamodb.AttributeValue{"PartitionKey": {S: aws.String("4")}},
				}, nil),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteTable(tableName),
			},
		},
	}

	grill.Run(t, tests)
}
