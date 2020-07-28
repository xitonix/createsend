package createsend

// ClientErrorCode client side error codes.
type ClientErrorCode int

const (
	// ErrCodeDataProcessing indicates that processing the input/output data has failed.
	ErrCodeDataProcessing ClientErrorCode = -1
	// ErrCodeNilHTTPClient the provided internal HTTP client is nil.
	ErrCodeNilHTTPClient ClientErrorCode = -2
	// ErrCodeAuthenticationNotSet neither API key nor Oauth token was provided.
	ErrCodeAuthenticationNotSet ClientErrorCode = -3
	// ErrCodeEmptyOAuthToken the provided Oauth token was empty.
	ErrCodeEmptyOAuthToken ClientErrorCode = -4
	// ErrCodeEmptyAPIKey the provided API key was empty.
	ErrCodeEmptyAPIKey ClientErrorCode = -5
	// ErrCodeEmptyURL the requested URL was empty.
	ErrCodeEmptyURL ClientErrorCode = -6
	// ErrCodeInvalidURL the requested UTL was invalid.
	ErrCodeInvalidURL ClientErrorCode = -7
	// ErrCodeInvalidJSON the provided JSON payload was invalid.
	ErrCodeInvalidJSON ClientErrorCode = -8
	// ErrCodeInvalidRequestBody the provided request was invalid.
	ErrCodeInvalidRequestBody ClientErrorCode = -9
)

// String returns the string representation of the error code.
func (c ClientErrorCode) String() string {
	switch c {
	case ErrCodeNilHTTPClient:
		return "The HTTP client cannot be nil"
	case ErrCodeAuthenticationNotSet:
		return "Either API key or OAuth authentication must be selected"
	case ErrCodeEmptyOAuthToken:
		return "the provided OAuth token was empty"
	case ErrCodeEmptyAPIKey:
		return "the provided API key was empty"
	case ErrCodeEmptyURL:
		return "the provided URL was empty"
	case ErrCodeInvalidURL:
		return "the provided URL is invalid"
	case ErrCodeInvalidJSON:
		return "invalid JSON data"
	case ErrCodeInvalidRequestBody:
		return "invalid request body"
	default:
		return "data processing error"
	}
}
