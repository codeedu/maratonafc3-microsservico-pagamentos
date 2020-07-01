package payment_providers

// This is the base struct when an error occurs during a request to the payment gateway
type ErrorProvider struct {
	Errors []Error `json:"errors"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// Detailed error
type Error struct {
	Message       string `json:"message"`
	ParameterName string `json:"parameter_name"`
	Type          string `json:"type"`
}


// Creates a new ErrorProvider
func NewErrorProvider() *ErrorProvider {
	return &ErrorProvider{}
}

// Creates a new Error struct to be used on an Error Provider
func NewError() *Error {
	return &Error{}
}