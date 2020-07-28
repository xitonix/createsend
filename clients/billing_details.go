package clients

// BillingMode billing mode.
type BillingMode int8

const (
	// MonthlyBilling monthly billing.
	MonthlyBilling BillingMode = iota + 1
	// PAYGBilling Pay as you go billing.
	PAYGBilling
)

// BillingDetails represents Clients' billing details.
type BillingDetails struct {
	// Mode billing mode.
	Mode BillingMode
	// Monthly is monthly billing details if the Client is on a monthly plan, nil otherwise.
	Monthly *MonthlyBillingDetails
	// PAYG is PAYG billing details of the Client is not on a monthly plan, nil otherwise.
	PAYG *PayAsYouGoBillingDetails
	// ClientPays is true if the client pays for itself.
	ClientPays bool
	// Currency the current billing currency.
	Currency string
}

// MonthlyBillingDetails represents monthly billing details if the Client is on a monthly plan.
type MonthlyBillingDetails struct {
	// Tier monthly billing tier.
	Tier string
	// Scheme the current scheme (eg. Unlimited, Basic, etc).
	Scheme string
	// Rate the current rate.
	Rate float64
	// MarkupPercentage the current markup percentage value.
	MarkupPercentage int64
	// Pending returns the plan pending for approval (if any).
	Pending *BillingDetails
}

// PayAsYouGoBillingDetails represents PAYG billing details if the client is on pay-as-you-go.
type PayAsYouGoBillingDetails struct {
	// CanPurchaseCredits is true if the Client is allowed to purchase credit.
	CanPurchaseCredits bool
	// Credits the number of Campaign Monitor credits.
	Credits int64
	// MarkupOnDesignSpamTest the current markup value on design and spam testing.
	MarkupOnDesignSpamTest float64
	// BaseRatePerRecipient the current rate per recipient.
	BaseRatePerRecipient float64
	// MarkupPerRecipient the current value of markup per recipient.
	MarkupPerRecipient float64
	// MarkupOnDelivery the current value of markup on delivery.
	MarkupOnDelivery float64
	// BaseDeliveryRate the current value of base delivery rate.
	BaseDeliveryRate float64
	// BaseDesignSpamTestRate the base rate for design & spam tests.
	BaseDesignSpamTestRate float64
}
