package application

import (
	"codeshop-payment/domain"
	"codeshop-payment/payment_providers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

func ProcessSubscription(transaction *domain.TransactionSubscriptionRequest) (domain.TransactionResponseInterface, error) {

	if transaction.Gateway.Name == "pagar.me" {
		pagarme := payment_providers.NewPagarme()
		pagarme.TransactionType = "subscription"
		pagarme.SubscriptionRequest = transaction
		pagarme.SubscriptionEndPoint = "https://api.pagar.me/1/subscriptions"
		err := pagarme.Process()

		if err != nil {
			return nil, err
		}

		if len(pagarme.Error.Errors) > 0 {
			return pagarme.Error, nil
		}

		clientResult, _ := GetClientByKey(transaction.SecretKey)
		client := domain.NewClient()
		client.ID = clientResult.ID

		pagarme.TransactionSubscriptionResponse.Client = client

		PersistTransaction(pagarme.TransactionSubscriptionResponse)
		return pagarme.TransactionSubscriptionResponse, nil
	}

	error := payment_providers.NewError()
	error.Message = "invalid gateway"

	return error, nil
}

func ProcessTransaction(transaction *domain.TransactionRequest) (domain.TransactionResponseInterface, error) {

	if transaction.Gateway.Name == "pagar.me" {
		pagarme := payment_providers.NewPagarme()
		pagarme.TransactionType = "transaction"
		pagarme.TransactionRequest = transaction
		pagarme.TransactionEndPoint = "https://api.pagar.me/1/transactions"
		err := pagarme.Process()

		if err != nil {
			return nil, err
		}

		if len(pagarme.Error.Errors) > 0 {
			return pagarme.Error, nil
		}

		clientResult, _ := GetClientByKey(transaction.SecretKey)
		client := domain.NewClient()
		client.ID = clientResult.ID

		pagarme.TransactionResponse.Client = client
		PersistTransaction(pagarme.TransactionResponse)

		return pagarme.TransactionResponse, nil
	}

	error := payment_providers.NewError()
	error.Message = "invalid gateway"

	return error, nil
}

func PersistTransaction(transaction domain.TransactionResponseInterface) error {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("region"))}))
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(transaction)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("transactions"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
