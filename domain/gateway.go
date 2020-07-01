package domain

// The Gateway is the payment gateway responsible to process a transaction. Ex: pagar.me
type Gateway struct {
	ID            string `json:"-"`
	Name          string `json:"name"`
	ApiKey        string `json:"api_key"`
	EncryptionKey string `json:"encryption_key"`
}
