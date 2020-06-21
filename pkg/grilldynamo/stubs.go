package grilldynamo

import (
	"bufio"
	"encoding/json"
	"os"

	"bitbucket.org/swigy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (grilldynamo *GrillDynamo) CreateTable(req *dynamodb.CreateTableInput) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := grilldynamo.dynamo.Client.CreateTable(req)
		return err
	})
}

// Reads a File containing json objects and puts them in DynamoDB.
// Each Line represents a new item
//
// Args -
//		tableName  - table to put data
//		filePath - absolute file path.
//
func (grilldynamo *GrillDynamo) SeedDataFromFile(tableName string, filePath string) grill.Stub {
	return grill.StubFunc(func() error {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			obj := scanner.Text()
			var input map[string]interface{}
			json.Unmarshal([]byte(obj), &input)
			item, err := dynamodbattribute.MarshalMap(&input)
			if err != nil {
				continue
			}
			_, err = grilldynamo.dynamo.Client.PutItem(&dynamodb.PutItemInput{
				Item:      item,
				TableName: aws.String(tableName),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
