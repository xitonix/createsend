package clients

// BasicDetails represents a client's basic details.
type BasicDetails struct {
	// Company company name.
	Company string `json:"CompanyName"`
	// Country country.
	Country string
	// Timezone client timezone (eg. "(GMT+10:00) Canberra, Melbourne, Sydney").
	Timezone string `json:"TimeZone"`
}
