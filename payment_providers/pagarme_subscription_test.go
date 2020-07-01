package payment_providers_test

import (
	"codeshop-payment/domain"
	"codeshop-payment/payment_providers"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var subscriptionRequest domain.TransactionSubscriptionRequest
var subscriptionRequestJson []byte
var pagarSubscriptionEndpoint string

func init() {
	pagarSubscriptionEndpoint = "/1/subscriptions"
	subscriptionRequest.APIKey = "key"
	subscriptionRequest.RemotePlanID = 123
	subscriptionRequest.PaymentMethod = "credit_card"
	subscriptionRequest.CardHash = "hash"
	subscriptionRequest.SoftDescriptor = "description"
	subscriptionRequest.PostbackURL = "http://localhost"

	customer := domain.CustomerSubscription{
		CustomerName:   "Wesley",
		CustomerEmail:  "wesley-test@email.com",
		DocumentNumber: "123.456.789-00",
	}

	subscriptionRequest.Customer = &customer
	subscriptionRequestJson, _ = json.Marshal(subscriptionRequest)
}

func TestPagarmeProcessApproved(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	transactionResponseExpected := &domain.TransactionSubscriptionResponse{
		ID:                   "",
		Client:               nil,
		Provider:             nil,
		ProcessType:          "subscription",
		RemoteSubscriptionID: 500030,
		Status:               "paid",
		CurrentTransaction: struct {
			RemoteTransactionID  int    `json:"id"`
			Amount               int    `json:"amount"`
			Installments         int    `json:"installments"`
			BoletoURL            string `json:"boleto_url"`
			BoletoBarcode        string `json:"boleto_barcode"`
			BoletoExpirationDate string `json:"boleto_expiration_date"`
		}{
			RemoteTransactionID:  9040653,
			Amount:               10000,
			Installments:         1,
			BoletoURL:            "",
			BoletoBarcode:        "",
			BoletoExpirationDate: "",
		},
		PaymentMethod:      "credit_card",
		CardBrand:          "visa",
		RemotePlanID:       486868,
		PostbackURL:        "http://google.com",
		CardLastDigits:     "1111",
		SoftDescriptor:     "",
		CurrentPeriodStart: "2020-06-22T21:41:01.468Z",
		CurrentPeriodSEnd:  "2021-06-22T21:41:01.468Z",
		RefuseReason:       "",
	}

	pagarme := payment_providers.NewPagarme()
	pagarme.ServiceAddress = srv.URL
	pagarme.TransactionType = "subscription"
	pagarme.SubscriptionRequest = &subscriptionRequest
	pagarme.SubscriptionEndPoint = pagarSubscriptionEndpoint
	pagarme.Process()
	pagarme.TransactionSubscriptionResponse.RemotePlanID = 486868

	transactionResponseExpected.ID = pagarme.TransactionSubscriptionResponse.ID
	transactionResponseExpected.Client = pagarme.TransactionSubscriptionResponse.Client
	transactionResponseExpected.Provider = pagarme.TransactionSubscriptionResponse.Provider
	require.Equal(t, transactionResponseExpected, pagarme.TransactionSubscriptionResponse)

}

func TestPagarmeProcessDeclined(t *testing.T) {
	srv := serverMockDeclined()
	defer srv.Close()

	error := payment_providers.ErrorProvider{
		Errors: []payment_providers.Error{
			{Type: "action_forbidden", ParameterName: "", Message: "Não foi possível realizar uma transação nesse cartão de crédito."},
		},
		Method: "post",
		URL: "/subscriptions?api_key=key",
	}

	pagarme := payment_providers.NewPagarme()
	pagarme.ServiceAddress = srv.URL
	pagarme.TransactionType = "subscription"
	pagarme.SubscriptionRequest = &subscriptionRequest
	pagarme.SubscriptionEndPoint = pagarSubscriptionEndpoint
	pagarme.Process()
	pagarme.TransactionSubscriptionResponse.RemotePlanID = 486868

	require.Equal(t, error, pagarme.Error)

}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(pagarSubscriptionEndpoint, processMock)
	srv := httptest.NewServer(handler)
	return srv
}

func processMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"object":"subscription","plan":{"object":"plan","id":486868,"amount":10000,"days":365,"name":"Plano anual","trial_days":0,"date_created":"2020-06-18T15:40:22.324Z","payment_methods":["boleto","credit_card"],"color":null,"charges":null,"installments":1,"invoice_reminder":null,"payment_deadline_charges_interval":1},"id":500030,"current_transaction":{"object":"transaction","status":"paid","refuse_reason":null,"status_reason":"acquirer","acquirer_response_code":"0000","acquirer_name":"pagarme","acquirer_id":"5eeb7b4347112365a5e8c7b2","authorization_code":"74125","soft_descriptor":null,"tid":9040653,"nsu":9040653,"date_created":"2020-06-22T21:41:01.497Z","date_updated":"2020-06-22T21:41:02.036Z","amount":10000,"authorized_amount":10000,"paid_amount":10000,"refunded_amount":0,"installments":1,"id":9040653,"cost":120,"card_holder_name":"Wes","card_last_digits":"1111","card_first_digits":"411111","card_brand":"visa","card_pin_mode":null,"card_magstripe_fallback":false,"cvm_pin":false,"postback_url":null,"payment_method":"credit_card","capture_method":"ecommerce","antifraud_score":null,"boleto_url":null,"boleto_barcode":null,"boleto_expiration_date":null,"referer":"api_key","ip":"216.239.181.161","subscription_id":500030,"metadata":{},"antifraud_metadata":{},"reference_key":null,"device":null,"local_transaction_id":null,"local_time":null,"fraud_covered":false,"fraud_reimbursed":null,"order_id":null,"risk_level":"very_low","receipt_url":null,"payment":null,"addition":null,"discount":null,"private_label":null},"postback_url":"http://google.com","payment_method":"credit_card","card_brand":"visa","card_last_digits":"1111","current_period_start":"2020-06-22T21:41:01.468Z","current_period_end":"2021-06-22T21:41:01.468Z","charges":0,"soft_descriptor":null,"status":"paid","date_created":"2020-06-22T21:41:02.026Z","date_updated":"2020-06-22T21:41:02.026Z","phone":null,"address":null,"customer":{"object":"customer","id":3330477,"external_id":null,"type":null,"country":null,"document_number":"31665663090","document_type":"cpf","name":"Wesley","email":"wesley-test@email.com","phone_numbers":null,"born_at":null,"birthday":null,"gender":null,"date_created":"2020-06-22T21:41:01.449Z","documents":[]},"card":{"object":"card","id":"card_ckbr0y0mh0b48my6d6gl1tj39","date_created":"2020-06-22T21:41:01.482Z","date_updated":"2020-06-22T21:41:02.019Z","brand":"visa","holder_name":"Wes","first_digits":"411111","last_digits":"1111","country":"UNITED STATES","fingerprint":"cj5bw4cio00000j23jx5l60cq","valid":true,"expiration_date":"1022"},"metadata":null,"fine":{},"interest":{},"settled_charges":null,"manage_token":"test_subscription_VVxd3bCo5oifIox5OECuxZm1nhNPjF","manage_url":"https://pagar.me/customers/#/subscriptions/500030?token=test_subscription_VVxd3bCo5oifIox5OECuxZm1nhNPjF"}`))
}

func serverMockDeclined() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(pagarSubscriptionEndpoint, processMockDeclined)
	srv := httptest.NewServer(handler)
	return srv
}

func processMockDeclined(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"errors":[{"type":"action_forbidden","parameter_name":null,"message":"Não foi possível realizar uma transação nesse cartão de crédito."}],"url":"/subscriptions?api_key=key","method":"post"}`))
}
