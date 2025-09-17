package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	Client    *dynamodb.DynamoDB
	TableName string
}

func NewDynamoDBClient() (*DynamoDBClient, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	tableName := os.Getenv("PRODUCTS_TABLE")
	if tableName == "" {
		tableName = "products-db"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := dynamodb.New(sess)

	return &DynamoDBClient{
		Client:    client,
		TableName: tableName,
	}, nil
}