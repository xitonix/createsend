package transactional

// Options represents Transactional API options.
type Options struct {
	clientID         string
	smartEmailStatus SmartEmailStatus
}

// Option represents a Transactional API option.
type Option func(options *Options)

// WithClientID sets the optional client ID.
func WithClientID(clientID string) Option {
	return func(options *Options) {
		options.clientID = clientID
	}
}

// ClientID returns the optional client ID.
func (o *Options) ClientID() string {
	return o.clientID
}

// WithSmartEmailStatus sets the optional smart email status.
func WithSmartEmailStatus(status SmartEmailStatus) Option {
	return func(options *Options) {
		options.smartEmailStatus = status
	}
}

// SmartEmailStatus returns the optional smart email status.
func (o *Options) SmartEmailStatus() SmartEmailStatus {
	return o.smartEmailStatus
}
