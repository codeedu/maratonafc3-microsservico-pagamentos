package payment_providers_test

import (
	"codeshop-payment/domain"
	"codeshop-payment/payment_providers"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

var transactionResponseJsonSubscriptionApproved string
var transactionResponseJsonSubscriptionError string
var transactionResponseJsonApproved string

func init() {
	transactionResponseJsonSubscriptionApproved = `{"object":"subscription","plan":{"object":"plan","id":486868,"amount":10000,"days":365,"name":"Plano anual","trial_days":0,"date_created":"2020-06-18T15:40:22.324Z","payment_methods":["boleto","credit_card"],"color":null,"charges":null,"installments":1,"invoice_reminder":null,"payment_deadline_charges_interval":1},"id":500030,"current_transaction":{"object":"transaction","status":"paid","refuse_reason":null,"status_reason":"acquirer","acquirer_response_code":"0000","acquirer_name":"pagarme","acquirer_id":"5eeb7b4347112365a5e8c7b2","authorization_code":"74125","soft_descriptor":null,"tid":9040653,"nsu":9040653,"date_created":"2020-06-22T21:41:01.497Z","date_updated":"2020-06-22T21:41:02.036Z","amount":10000,"authorized_amount":10000,"paid_amount":10000,"refunded_amount":0,"installments":1,"id":9040653,"cost":120,"card_holder_name":"Wes","card_last_digits":"1111","card_first_digits":"411111","card_brand":"visa","card_pin_mode":null,"card_magstripe_fallback":false,"cvm_pin":false,"postback_url":null,"payment_method":"credit_card","capture_method":"ecommerce","antifraud_score":null,"boleto_url":null,"boleto_barcode":null,"boleto_expiration_date":null,"referer":"api_key","ip":"216.239.181.161","subscription_id":500030,"metadata":{},"antifraud_metadata":{},"reference_key":null,"device":null,"local_transaction_id":null,"local_time":null,"fraud_covered":false,"fraud_reimbursed":null,"order_id":null,"risk_level":"very_low","receipt_url":null,"payment":null,"addition":null,"discount":null,"private_label":null},"postback_url":"http://google.com","payment_method":"credit_card","card_brand":"visa","card_last_digits":"1111","current_period_start":"2020-06-22T21:41:01.468Z","current_period_end":"2021-06-22T21:41:01.468Z","charges":0,"soft_descriptor":null,"status":"paid","date_created":"2020-06-22T21:41:02.026Z","date_updated":"2020-06-22T21:41:02.026Z","phone":null,"address":null,"customer":{"object":"customer","id":3330477,"external_id":null,"type":null,"country":null,"document_number":"31665663090","document_type":"cpf","name":"Wesley","email":"wesley-test@email.com","phone_numbers":null,"born_at":null,"birthday":null,"gender":null,"date_created":"2020-06-22T21:41:01.449Z","documents":[]},"card":{"object":"card","id":"card_ckbr0y0mh0b48my6d6gl1tj39","date_created":"2020-06-22T21:41:01.482Z","date_updated":"2020-06-22T21:41:02.019Z","brand":"visa","holder_name":"Wes","first_digits":"411111","last_digits":"1111","country":"UNITED STATES","fingerprint":"cj5bw4cio00000j23jx5l60cq","valid":true,"expiration_date":"1022"},"metadata":null,"fine":{},"interest":{},"settled_charges":null,"manage_token":"test_subscription_VVxd3bCo5oifIox5OECuxZm1nhNPjF","manage_url":"https://pagar.me/customers/#/subscriptions/500030?token=test_subscription_VVxd3bCo5oifIox5OECuxZm1nhNPjF"}`
	transactionResponseJsonSubscriptionError = `{"errors":[{"type":"action_forbidden","parameter_name":null,"message":"Não foi possível realizar uma transação nesse cartão de crédito."}],"url":"/subscriptions?api_key=key","method":"post"}`
	transactionResponseJsonApproved = `{"object":"transaction","status":"paid","refuse_reason":null,"status_reason":"acquirer","acquirer_response_code":"0000","acquirer_name":"pagarme","acquirer_id":"5eeb7b4347112365a5e8c7b2","authorization_code":"439676","soft_descriptor":null,"tid":9079267,"nsu":9079267,"date_created":"2020-06-29T01:07:43.170Z","date_updated":"2020-06-29T01:07:43.569Z","amount":10000,"authorized_amount":10000,"paid_amount":10000,"refunded_amount":0,"installments":1,"id":9079267,"cost":50,"card_holder_name":"Morpheus Fishburne","card_last_digits":"1111","card_first_digits":"411111","card_brand":"visa","card_pin_mode":null,"card_magstripe_fallback":false,"cvm_pin":false,"postback_url":null,"payment_method":"credit_card","capture_method":"ecommerce","antifraud_score":null,"boleto_url":null,"boleto_barcode":null,"boleto_expiration_date":null,"referer":"api_key","ip":"216.239.181.161","subscription_id":null,"phone":null,"address":null,"customer":{"object":"customer","id":3348969,"external_id":"#3311","type":"individual","country":"br","document_number":null,"document_type":"cpf","name":"Morpheus Fishburne","email":"mopheus@nabucodonozor.com","phone_numbers":["+5511999998888"],"born_at":null,"birthday":null,"gender":null,"date_created":"2020-06-29T01:07:43.111Z","documents":[{"object":"document","id":"doc_ckbzsyxsj0a6oj86dhgnx53kh","type":"cpf","number":"30621143049"}]},"billing":null,"shipping":null,"items":[],"card":{"object":"card","id":"card_ckbzsyxth0a6pj86d4ndteluf","date_created":"2020-06-29T01:07:43.158Z","date_updated":"2020-06-29T01:07:43.640Z","brand":"visa","holder_name":"Morpheus Fishburne","first_digits":"411111","last_digits":"1111","country":"UNITED STATES","fingerprint":"cj5bw4cio00000j23jx5l60cq","valid":true,"expiration_date":"0922"},"split_rules":null,"metadata":{},"antifraud_metadata":{},"reference_key":null,"device":null,"local_transaction_id":null,"local_time":null,"fraud_covered":false,"fraud_reimbursed":null,"order_id":null,"risk_level":"unknown","receipt_url":null,"payment":null,"addition":null,"discount":null,"private_label":null}`
}

func TestConvertJsonToSubscriptionTransactionApprovedResponse(t *testing.T) {

	uuid := uuid.NewV4().String()
	transactionResponse := domain.TransactionSubscriptionResponse{
		ID:                   uuid,
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

	var transactionResponseToTest domain.TransactionSubscriptionResponse
	json.Unmarshal([]byte(transactionResponseJsonSubscriptionApproved), &transactionResponseToTest)
	transactionResponseToTest.RemotePlanID = 486868
	transactionResponseToTest.ID = uuid

	require.Equal(t, transactionResponse, transactionResponseToTest)
}

func TestConvertsJsonToTransactionSubscriptionDeclineddResponse(t *testing.T) {

	error := payment_providers.ErrorProvider{
		Errors: []payment_providers.Error{
			{Type: "action_forbidden", ParameterName: "", Message: "Não foi possível realizar uma transação nesse cartão de crédito."},
		},
		Method: "post",
		URL:    "/subscriptions?api_key=key",
	}

	var transactionResponseToTest payment_providers.ErrorProvider
	json.Unmarshal([]byte(transactionResponseJsonSubscriptionError), &transactionResponseToTest)

	require.Equal(t, error, transactionResponseToTest)
}

func TestConvertJsonToTransactionApprovedResponse(t *testing.T) {

	uuid := uuid.NewV4().String()

	var phoneNumbers []string
	phoneNumbers = append(phoneNumbers, "+5511999998888")

	transactionResponse := domain.TransactionResponse{
		ID:                    uuid,
		Client:                nil,
		Provider:              nil,
		ProcessType:           "transaction",
		Status:                "paid",
		RefuseReason:          "",
		StatusReason:          "acquirer",
		AcquirerResponseCode:  "0000",
		AcquirerName:          "pagarme",
		AcquirerId:            "5eeb7b4347112365a5e8c7b2",
		AuthorizationCode:     "439676",
		SoftDescriptor:        "",
		TID:                   9079267,
		NSU:                   9079267,
		CreatedAt:             "2020-06-29T01:07:43.170Z",
		UpdatedAt:             "2020-06-29T01:07:43.569Z",
		Amount:                10000,
		AuthorizedAmount:      10000,
		PaidAmount:            10000,
		RefundedAmount:        0,
		Installments:          1,
		RemoteTransactionID:   9079267,
		Cost:                  50,
		CardHolderName:        "Morpheus Fishburne",
		CardLastDigits:        "1111",
		CardFirstDigits:       "411111",
		CardBrand:             "visa",
		CardPinMode:           "",
		CardMagstripeFallback: false,
		CvmPin:                false,
		PostbackURL:           "",
		PaymentMethod:         "credit_card",
		CaptureMethod:         "ecommerce",
		AntifraudScore:        "",
		BoletoURL:             "",
		BoletoBarcode:         "",
		BoletoExpirationDate:  "",
		Referer:               "api_key",
		IP:                    "216.239.181.161",
		SubscriptionId:        "",
		Phone:                 "",
		Address:               "",
		Customer: domain.Customer{
			ID:               "",
			Object:           "customer",
			RemoteCustomerID: 3348969,
			ExternalId:       "#3311",
			CustomerType:     "individual",
			Country:          "br",
			DocumentNumber:   "",
			DocumentType:     "cpf",
			Name:             "Morpheus Fishburne",
			Email:            "mopheus@nabucodonozor.com",
			PhoneNumbers: phoneNumbers,
			BornAt:      "",
			Birthday:    "",
			Gender:      "",
			DateCreated: "2020-06-29T01:07:43.111Z",
			Documents: []domain.CustomerDocuments{
				{

					Object:       "document",
					ID:           "doc_ckbzsyxsj0a6oj86dhgnx53kh",
					DocumentType: "cpf",
					Number:       "30621143049",
				},
			},
		},
		Billing:  "",
		Shipping: "",
		Items: []string{},
		Card: domain.CardResponseTransaction{
			Object:         "card",
			ID:             "card_ckbzsyxth0a6pj86d4ndteluf",
			DateCreated:    "2020-06-29T01:07:43.158Z",
			DateUpdated:    "2020-06-29T01:07:43.640Z",
			Brand:          "visa",
			HolderName:     "Morpheus Fishburne",
			FirstDigits:    "411111",
			LastDigits:     "1111",
			Country:        "UNITED STATES",
			Fingerprint:    "cj5bw4cio00000j23jx5l60cq",
			Valid:          true,
			ExpirationDate: "0922",
		},
		SplitRules:          "",
		Metadata:            struct{}{},
		AntifraudMetadata:   struct{}{},
		ReferenceKey:        "",
		Device:              "",
		LocalTransaction_id: "",
		LocalTime:           "",
		FraudCovered:        false,
		FraudReimbursed:     "",
		OrderId:             "",
		RiskLevel:           "unknown",
		ReceiptUrl:          "",
		Payment:             "",
		Addition:            "",
		Discount:            "",
		PrivateLabel:        "",
	}

	var transactionResponseToTest domain.TransactionResponse
	json.Unmarshal([]byte(transactionResponseJsonApproved), &transactionResponseToTest)
	transactionResponseToTest.ID = uuid
	require.Equal(t, transactionResponse, transactionResponseToTest)
}
