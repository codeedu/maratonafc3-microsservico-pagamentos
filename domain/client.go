package domain

// Client means every system which is going to use this microservice
// In order to make sure that this application can recognize which client is making a request, a key is required
// Also the goal of the key is to keep this microservice safe by avoiding others systems making not authorized requests
type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SecretKey  string `json:"secret_key"`
}

// The NewClient function is responsible to create an empty Client struct
func NewClient() *Client {
	return &Client{}
}
