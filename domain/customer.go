package domain

// The CustomerInterface is an empty interface that allow us to use either Customer or CustomerSubscription structs as params
type CustomerInterface interface{}

// This struct is responsible for handle customer information for regular transactions
type Customer struct {
	ID               string              `json:"-"`
	Object           string              `json:"object"`
	RemoteCustomerID int                 `json:"id"`
	ExternalId       string              `json:"external_id"`
	CustomerType     string              `json:"type"`
	Country          string              `json:"country"`
	DocumentNumber   string              `json:"document_number"`
	DocumentType     string              `json:"document_type"`
	Name             string              `json:"name"`
	Email            string              `json:"email"`
	PhoneNumbers     []string            `json:"phone_numbers"`
	BornAt           string              `json:"born_at"`
	Birthday         string              `json:"birthday"`
	Gender           string              `json:"gender"`
	DateCreated      string           `json:"date_created"`
	Documents        []CustomerDocuments `json:"documents"`
}

// CustomerDocuments aggregates all documents which are required to process a transaction
type CustomerDocuments struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	DocumentType string `json:"type"`
	Number       string `json:"number"`
}

// This struct is responsible for handle customer information for transactions used for subscriptions only
type CustomerSubscription struct {
	ID             string `json:"-"`
	CustomerName   string `json:"name"`
	CustomerEmail  string `json:"email"`
	DocumentNumber string `json:"document_number"`
}
