package payment_providers

import (
	"bytes"
	"codeshop-payment/domain"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
)

// This is the concrete implementation of pagar.me payment gateway
// Its supported two types of requests: "transaction" and "subscription"
type Pagarme struct {
	SubscriptionRequest             *domain.TransactionSubscriptionRequest
	TransactionRequest              *domain.TransactionRequest
	TransactionSubscriptionResponse *domain.TransactionSubscriptionResponse
	TransactionResponse             *domain.TransactionResponse
	Error                           ErrorProvider
	SubscriptionEndPoint            string
	TransactionEndPoint             string
	TransactionType                 string // subscription or transaction
	ServiceAddress                  string
}

// Created a new instance of Pagarme
func NewPagarme() *Pagarme {
	return &Pagarme{}
}

func (p *Pagarme) prepareRequest() ([]byte, error) {

	var requestJson []byte
	var err error

	if p.TransactionType == "subscription" {
		requestJson, err = json.Marshal(p.SubscriptionRequest)
	} else {
		requestJson, err = json.Marshal(p.TransactionRequest)
	}

	if err != nil {
		return nil, err
	}

	return requestJson, nil
}

// Process the transaction regardless if its a regular one or a subscription
func (p *Pagarme) Process() error {

	dataRequestJson, _ := p.prepareRequest()

	var endpoint string

	if p.TransactionType == "subscription" {
		endpoint = p.SubscriptionEndPoint
	} else {
		endpoint = p.TransactionEndPoint
	}

	req, err := http.NewRequest("POST", p.ServiceAddress+endpoint, bytes.NewBuffer(dataRequestJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	transactionId := uuid.NewV4().String()

	if p.TransactionType == "subscription" {
		json.Unmarshal(body, &p.TransactionSubscriptionResponse)
		p.TransactionSubscriptionResponse.ID = transactionId
		p.TransactionSubscriptionResponse.Provider = &domain.Gateway{Name: "pagar.me"}

		if p.TransactionSubscriptionResponse.CurrentTransaction.RemoteTransactionID == 0 {
			json.Unmarshal(body, &p.Error)
		}

	} else {
		json.Unmarshal(body, &p.TransactionResponse)
		p.TransactionResponse.ID = transactionId

		if p.TransactionResponse.RemoteTransactionID == 0 {
			json.Unmarshal(body, &p.Error)
		}
	}

	return nil
}
