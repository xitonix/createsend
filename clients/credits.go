package clients

// CreditTransferRequest credit transfer request.
type CreditTransferRequest struct {
	// Credits number of credits to transfer.
	Credits int
	// CanUseMyCreditsWhenTheyRunOut if true the client will be able to continue sending emails using your credits or payment details once they run out of credits.
	//
	// Otherwise the client will not be able to continue sending until you allocate more credits to them.
	CanUseMyCreditsWhenTheyRunOut bool
}

// CreditTransferResult credit transfer result.
type CreditTransferResult struct {
	// AccountCredits the source/target account's credits.
	AccountCredits int
	// ClientCredits the source/target client's credits.
	ClientCredits int
}
