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

var transactionRequest domain.TransactionRequest
var transactionRequestJson []byte
var pagarmeTransactionEndpoint string

func init() {
	pagarmeTransactionEndpoint = "/1/transactions"
	transactionRequest.APIKey = "key"
	transactionRequest.PaymentMethod = "credit_card"
	transactionRequest.CardHash = "hash"
	transactionRequest.SoftDescriptor = "description"
	transactionRequest.PostbackURL = "http://localhost"
	CustomerDocuments := []domain.CustomerDocuments{
		{
			DocumentType: "cpf",
			Number:       "30621143049",
		},
	}

	customer := domain.Customer{
		ExternalId:   "#3311",
		Name:         "Morpheus Fishburne",
		CustomerType: "individual",
		Country:      "br",
		Email:        "mopheus@nabucodonozor.com",
		Documents:    CustomerDocuments,
		PhoneNumbers: []string{
			"+5511999998888",
		},
	}

	transactionRequest.Customer = customer
	transactionRequestJson, _ = json.Marshal(transactionRequest)
}

func TestPagarmeTransactionProcessApproved(t *testing.T) {
	srv := serverTransactionMock()
	defer srv.Close()

	transactionResponseExpected := domain.TransactionResponse{
		ID:                    "",
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
			PhoneNumbers: []string{
				"+5511999998888",
			},
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
		Items:    []string{},
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

	pagarme := payment_providers.NewPagarme()
	pagarme.ServiceAddress = srv.URL
	pagarme.TransactionType = "transaction"
	pagarme.TransactionRequest = &transactionRequest
	pagarme.TransactionEndPoint = pagarmeTransactionEndpoint
	pagarme.Process()

	transactionResponseExpected.ID = pagarme.TransactionResponse.ID
	require.Equal(t, &transactionResponseExpected, pagarme.TransactionResponse)
}

func serverTransactionMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(pagarmeTransactionEndpoint, processTransactionMock)
	srv := httptest.NewServer(handler)
	return srv
}

func processTransactionMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"object":"transaction","status":"paid","refuse_reason":null,"status_reason":"acquirer","acquirer_response_code":"0000","acquirer_name":"pagarme","acquirer_id":"5eeb7b4347112365a5e8c7b2","authorization_code":"439676","soft_descriptor":null,"tid":9079267,"nsu":9079267,"date_created":"2020-06-29T01:07:43.170Z","date_updated":"2020-06-29T01:07:43.569Z","amount":10000,"authorized_amount":10000,"paid_amount":10000,"refunded_amount":0,"installments":1,"id":9079267,"cost":50,"card_holder_name":"Morpheus Fishburne","card_last_digits":"1111","card_first_digits":"411111","card_brand":"visa","card_pin_mode":null,"card_magstripe_fallback":false,"cvm_pin":false,"postback_url":null,"payment_method":"credit_card","capture_method":"ecommerce","antifraud_score":null,"boleto_url":null,"boleto_barcode":null,"boleto_expiration_date":null,"referer":"api_key","ip":"216.239.181.161","subscription_id":null,"phone":null,"address":null,"customer":{"object":"customer","id":3348969,"external_id":"#3311","type":"individual","country":"br","document_number":null,"document_type":"cpf","name":"Morpheus Fishburne","email":"mopheus@nabucodonozor.com","phone_numbers":["+5511999998888"],"born_at":null,"birthday":null,"gender":null,"date_created":"2020-06-29T01:07:43.111Z","documents":[{"object":"document","id":"doc_ckbzsyxsj0a6oj86dhgnx53kh","type":"cpf","number":"30621143049"}]},"billing":null,"shipping":null,"items":[],"card":{"object":"card","id":"card_ckbzsyxth0a6pj86d4ndteluf","date_created":"2020-06-29T01:07:43.158Z","date_updated":"2020-06-29T01:07:43.640Z","brand":"visa","holder_name":"Morpheus Fishburne","first_digits":"411111","last_digits":"1111","country":"UNITED STATES","fingerprint":"cj5bw4cio00000j23jx5l60cq","valid":true,"expiration_date":"0922"},"split_rules":null,"metadata":{},"antifraud_metadata":{},"reference_key":null,"device":null,"local_transaction_id":null,"local_time":null,"fraud_covered":false,"fraud_reimbursed":null,"order_id":null,"risk_level":"unknown","receipt_url":null,"payment":null,"addition":null,"discount":null,"private_label":null}`))
}

func serverMockTransactionDeclined() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(pagarSubscriptionEndpoint, processTransactionMockDeclined)
	srv := httptest.NewServer(handler)
	return srv
}

func processTransactionMockDeclined(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"errors":[{"type":"action_forbidden","parameter_name":null,"message":"Não foi possível realizar uma transação nesse cartão de crédito."}],"url":"/subscriptions?api_key=key","method":"post"}`))
}
