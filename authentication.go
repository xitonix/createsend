package createsend

import "net/http"

type authenticationMethod int

const (
	undefinedAuthentication authenticationMethod = iota
	apiKeyAuthentication
	oAuthAuthentication
)

type authentication struct {
	token  string
	method authenticationMethod
}

func (a *authentication) apply(request *http.Request) {
	if request == nil {
		return
	}
	switch a.method {
	case apiKeyAuthentication:
		request.SetBasicAuth(a.token, a.token)
	case oAuthAuthentication:
		request.Header.Set(authenticationHeaderKey, "Bearer "+a.token)
	}
}

func (a *authentication) validate() error {
	if a == nil || a.method == undefinedAuthentication {
		return newClientError(ErrCodeAuthenticationNotSet)
	}
	if len(a.token) == 0 {
		switch a.method {
		case apiKeyAuthentication:
			return newClientError(ErrCodeEmptyAPIKey)
		case oAuthAuthentication:
			return newClientError(ErrCodeEmptyOAuthToken)
		}
	}
	return nil
}
