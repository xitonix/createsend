package campaigns

// EmailClientUsage represents the information regarding the type of devices that were used to open emails
type EmailClientUsage struct {
	// Client represents the type of device
	Client string
	// Version represents the version of the device
	Version string
	// Percentage represents the percentage of the email opens for this specific device
	Percentage float32
	// Subscribers represents the total amount of subscribers that opened the email using this specific device
	Subscribers int
}
