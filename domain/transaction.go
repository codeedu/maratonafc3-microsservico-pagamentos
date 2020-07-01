package domain

import (
	"time"
)

// Interface which all payment providers have to implement
type PaymentProvider interface {
	Process() error
}

// TransactionRequest is responsible to handle data received from a regular transaction request
type TransactionRequest struct {
	ID                   string    `json:"-"`
	SecretKey            string    `json:"secret_key"`
	Client               *Client   `json:"-"`
	Gateway              Gateway   `json:"gateway"`
	APIKey               string    `json:"api_key"`
	TransactionClientID  string    `json:"transaction_client_id"`
	PaymentMethod        string    `json:"payment_method"`
	Amount               int       `json:"amount"`
	CardHash             string    `json:"card_hash"`
	PostbackURL          string    `json:"postback_url"`
	async                bool      `json:"async"`
	Installments         int       `json:"installments"`
	BoletoExpirationDate time.Time `json:"boleto_expiration_date"`
	SoftDescriptor       string    `json:"soft_descriptor"`
	Capture              bool      `json:"capture"`
	BoletoInstructions   string    `json:"boleto_instructions"`
	Customer             Customer  `json:"customer"`
}

// TransactionSubscriptionRequest is responsible to handle data received from a subscription request
type TransactionSubscriptionRequest struct {
	ID                    string                `json:"-"`
	SecretKey             string                `json:"secret_key"`
	Client                *Client               `json:"-"`
	Gateway               Gateway               `json:"gateway"`
	TransactionProviderID string                `json:"-"`
	APIKey                string                `json:"api_key"`
	RemotePlanID          int                   `json:"plan_id"`
	PaymentMethod         string                `json:"payment_method"`
	CardHash              string                `json:"card_hash"`
	SoftDescriptor        string                `json:"soft_descriptor"`
	PostbackURL           string                `json:"postback_url"`
	Customer              *CustomerSubscription `json:"customer"`
}

// This interface allow us to return different types of responses such as regular transactions and subscription responses
type TransactionResponseInterface interface {
}

// This is the response of a subscription transaction following the pagar.me's output
type TransactionSubscriptionResponse struct {
	ID                   string `json:"transaction_id"`
	Client               *Client
	Provider             *Gateway
	ProcessType          string `json:"object"`
	RemoteSubscriptionID int    `json:"id"`
	Status               string `json:"status"`
	CurrentTransaction   struct {
		RemoteTransactionID  int    `json:"id"`
		Amount               int    `json:"amount"`
		Installments         int    `json:"installments"`
		BoletoURL            string `json:"boleto_url"`
		BoletoBarcode        string `json:"boleto_barcode"`
		BoletoExpirationDate string `json:"boleto_expiration_date"`
	} `json:"current_transaction"`
	PaymentMethod      string    `json:"payment_method"`
	CardBrand          string    `json:"card_brand"`
	RemotePlanID       int       `json:"remote_plan_id"`
	PostbackURL        string    `json:"postback_url"`
	CardLastDigits     string    `json:"card_last_digits"`
	SoftDescriptor     string    `json:"soft_descriptor"`
	CurrentPeriodStart string    `json:"current_period_start"`
	CurrentPeriodSEnd  string    `json:"current_period_end"`
	RefuseReason       string    `json:"refuse_reason"`
	CreatedAt          time.Time `json:"date_created"`
	UpdatedAt          time.Time `json:"date_created"`
}

// This is the response of a REGULAR transaction following the pagar.me's output
type TransactionResponse struct {
	ID                    string `json:"transaction_id"`
	Client                *Client
	Provider              *Gateway
	ProcessType           string   `json:"object"`
	Status                string   `json:"status"`
	RefuseReason          string   `json:"refuse_reason"`
	StatusReason          string   `json:"status_reason"`
	AcquirerResponseCode  string   `json:"acquirer_response_code"`
	AcquirerName          string   `json:"acquirer_name"`
	AcquirerId            string   `json:"acquirer_id"`
	AuthorizationCode     string   `json:"authorization_code"`
	SoftDescriptor        string   `json:"soft_descriptor"`
	TID                   int      `json:"tid"`
	NSU                   int      `json:"nsu"`
	CreatedAt             string   `json:"date_created"`
	UpdatedAt             string   `json:"date_updated"`
	Amount                int      `json:"amount"`
	AuthorizedAmount      int      `json:"authorized_amount"`
	PaidAmount            int      `json:"paid_amount"`
	RefundedAmount        int      `json:"refunded_amount"`
	Installments          int      `json:"installments"`
	RemoteTransactionID   int      `json:"id"`
	Cost                  int      `json:"cost"`
	CardHolderName        string   `json:"card_holder_name"`
	CardLastDigits        string   `json:"card_last_digits"`
	CardFirstDigits       string   `json:"card_first_digits"`
	CardBrand             string   `json:"card_brand"`
	CardPinMode           string   `json:"card_pin_mode"`
	CardMagstripeFallback bool     `json:"card_magstripe_fallback"`
	CvmPin                bool     `json:"cvm_pin"`
	PostbackURL           string   `json:"postback_url"`
	PaymentMethod         string   `json:"payment_method"`
	CaptureMethod         string   `json:"capture_method"`
	AntifraudScore        string   `json:"antifraud_score"`
	BoletoURL             string   `json:"boleto_url"`
	BoletoBarcode         string   `json:"boleto_barcode"`
	BoletoExpirationDate  string   `json:"boleto_expiration_date"`
	Referer               string   `json:"referer"`
	IP                    string   `json:"ip"`
	SubscriptionId        string   `json:"subscription_id"`
	Phone                 string   `json:"phone"`
	Address               string   `json:"address"`
	Customer              Customer `json:"customer"`
	Billing               string   `json:"billing"`
	Shipping              string   `json:"shipping"`
	Items                 []string `json:"items"`
	Card       CardResponseTransaction `json:"card"`
	SplitRules string                  `json:"split_rules"`
	Metadata   struct {
	} `json:"metadata"`
	AntifraudMetadata struct {
	} `json:"antifraud_metadata"`
	ReferenceKey        string `json:"reference_key"`
	Device              string `json:"device"`
	LocalTransaction_id string `json:"local_transaction_id"`
	LocalTime           string `json:"local_time"`
	FraudCovered        bool `json:"fraud_covered"`
	FraudReimbursed     string `json:"fraud_reimbursed"`
	OrderId             string `json:"order_id"`
	RiskLevel           string `json:"risk_level"`
	ReceiptUrl          string `json:"receipt_url"`
	Payment             string `json:"payment"`
	Addition            string `json:"addition"`
	Discount            string `json:"discount"`
	PrivateLabel        string `json:"private_label"`
}

// Returns the credit card information in a transaction type response
type CardResponseTransaction struct {
	Object         string `json:"object"`
	ID             string `json:"id"`
	DateCreated    string `json:"date_created"`
	DateUpdated    string `json:"date_updated"`
	Brand          string `json:"brand"`
	HolderName     string `json:"holder_name"`
	FirstDigits    string `json:"first_digits"`
	LastDigits     string `json:"last_digits"`
	Country        string `json:"country"`
	Fingerprint    string `json:"fingerprint"`
	Valid          bool   `json:"valid"`
	ExpirationDate string `json:"expiration_date"`
}

// Creates a new regular transaction request
func NewTransactionRequest() *TransactionRequest {
	return &TransactionRequest{}
}

// Creates a new subscription transaction request
func NewTransactionSubscriptionRequest() *TransactionSubscriptionRequest {
	return &TransactionSubscriptionRequest{}
}

// Creates a new regular transaction request
func NewTransactionResponse() *TransactionResponse {
	return &TransactionResponse{}
}
