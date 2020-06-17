package clients

type ClientDetails struct {
	ApiKey              string
	Id                  string
	Company             string
	Country             string
	Timezone            string
	PrimaryContactName  string
	PrimaryContactEmail string
	Billing             BillingDetails
}
