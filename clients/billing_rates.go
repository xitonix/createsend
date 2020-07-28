package clients

// PAYGRates represents PAYG billing rates.
type PAYGRates struct {
	// Currency the billing currency.
	Currency string
	// CanPurchaseCredits is true if the Client is allowed to purchase credit.
	CanPurchaseCredits bool
	// ClientPays is true if the client pays for itself.
	ClientPays bool
	// MarkupPercentage markup percentage value.
	MarkupPercentage int64
	// MarkupOnDelivery markup on delivery.
	MarkupOnDelivery float64
	// MarkupPerRecipient markup per recipient.
	MarkupPerRecipient float64
	// MarkupOnDesignSpamTest markup value on design and spam testing.
	MarkupOnDesignSpamTest float64
}

// MonthlyRates represents monthly billing rates.
type MonthlyRates struct {
	// Currency the billing currency.
	Currency string
	// ClientPays is true if the client pays for itself.
	ClientPays bool
	// MarkupPercentage markup percentage value.
	MarkupPercentage int64
	// Scheme the billing scheme (eg. Unlimited, Basic, etc).
	Scheme string `json:"MonthlyScheme"`
}
