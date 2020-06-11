package createsend

type ClientErrorCode int

const (
	ErrCodeUnknown              ClientErrorCode = -1
	ErrCodeNilHTTPClient        ClientErrorCode = -2
	ErrCodeAuthenticationNotSet ClientErrorCode = -3
	ErrCodeEmptyOAuthToken      ClientErrorCode = -4
	ErrCodeEmptyAPIKey          ClientErrorCode = -5
	ErrCodeEmptyURL             ClientErrorCode = -6
	ErrCodeInvalidURL           ClientErrorCode = -7
	ErrCodeInvalidJson          ClientErrorCode = -8
	ErrCodeInvalidRequestBody   ClientErrorCode = -9
)

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
	case ErrCodeInvalidJson:
		return "invalid json data"
	case ErrCodeInvalidRequestBody:
		return "invalid request body"
	default:
		return "Unknown client error"
	}
}
