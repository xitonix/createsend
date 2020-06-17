package clients

type BillingDetails struct {
	CurrentTier        string
	CurrentMonthlyRate float64
	MarkupPercentage   float64
	ClientPays         bool
	MonthlyScheme      string
	Currency           string
}
