package main

import (
	"codeshop-payment/application"
	"codeshop-payment/domain"
	"codeshop-payment/payment_providers"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func SubscriptionHandler(request events.APIGatewayProxyRequest) (Response, error) {

	transaction := domain.NewTransactionSubscriptionRequest()
	json.Unmarshal([]byte(request.Body), &transaction)

	_, err := application.GetClientByKey(transaction.SecretKey)

	if err != nil {
		return returnResponse(err.Error()), nil
	}

	transactionResponse, err := application.ProcessSubscription(transaction)

	if err != nil {
		return returnResponse(err.Error()), nil
	}

	transactionJson, _ := json.Marshal(transactionResponse)

	return returnResponse(string(transactionJson)), nil
}

func process(transaction *domain.TransactionSubscriptionRequest) (Response, error) {

	if transaction.Gateway.Name == "pagar.me" {
		pagarme := payment_providers.NewPagarme()
		pagarme.TransactionType = "subscription"
		pagarme.SubscriptionRequest = transaction
		pagarme.SubscriptionEndPoint = "https://api.pagar.me/1/subscriptions"
		err := pagarme.Process()

		if err != nil {
			return returnResponse(err.Error()), nil
		}

		if len(pagarme.Error.Errors) > 0 {
			transactionJson, _ := json.Marshal(pagarme.Error)
			return returnResponse(string(transactionJson)), nil
		}

		transactionJson, _ := json.Marshal(pagarme.TransactionResponse)
		return returnResponse(string(transactionJson)), nil
	}

	return returnResponse(string(`{"error":"invalid gateway"}`)), nil
}

func main() {
	lambda.Start(SubscriptionHandler)
}

func returnResponse(body string) Response {
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            body,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp
}
