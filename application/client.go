package application

import (
	"codeshop-payment/domain"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"os"
)

func GetClientByKey(key string) (*domain.Client, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("region"))}))
	svc := dynamodb.New(sess)

	filt := expression.Name("secret_key").Equal(expression.Value(key))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		FilterExpression: expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:        aws.String("clients"),
		IndexName:        aws.String("key-index"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		return nil, err
	}

	if len(result.Items) > 0 {
		item := domain.NewClient()
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
		return item, nil
	}

	return nil, errors.New("Item not found")
}
