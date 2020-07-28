package internal

import (
	"strings"

	"github.com/xitonix/createsend/clients"
)

// BillingDetails represents a raw billing details type.
type BillingDetails struct {
	// Monthly plans
	CurrentTier        string
	MonthlyScheme      string
	CurrentMonthlyRate float64
	MarkupPercentage   int64

	// Pay as you go
	CanPurchaseCredits     bool
	Credits                int64
	MarkupOnDesignSpamTest float64
	BaseRatePerRecipient   float64
	MarkupPerRecipient     float64
	MarkupOnDelivery       float64
	BaseDeliveryRate       float64
	BaseDesignSpamTestRate float64

	// Common
	ClientPays bool
	Currency   string
}

// ToClientBillingDetails converts the raw model to a new createsend model.
func (b *BillingDetails) ToClientBillingDetails(pending *BillingDetails) *clients.BillingDetails {
	if b == nil {
		return nil
	}
	if len(strings.TrimSpace(b.CurrentTier)) > 0 {
		bd := &clients.BillingDetails{
			Mode: clients.MonthlyBilling,
			Monthly: &clients.MonthlyBillingDetails{
				Tier:             b.CurrentTier,
				Scheme:           b.MonthlyScheme,
				Rate:             b.CurrentMonthlyRate,
				MarkupPercentage: b.MarkupPercentage,
			},
			ClientPays: b.ClientPays,
			Currency:   b.Currency,
		}
		if pending != nil {
			bd.Monthly.Pending = pending.ToClientBillingDetails(nil)
		}
		return bd
	}

	return &clients.BillingDetails{
		Mode: clients.PAYGBilling,
		PAYG: &clients.PayAsYouGoBillingDetails{
			CanPurchaseCredits:     b.CanPurchaseCredits,
			Credits:                b.Credits,
			MarkupOnDesignSpamTest: b.MarkupOnDesignSpamTest,
			BaseRatePerRecipient:   b.BaseRatePerRecipient,
			MarkupPerRecipient:     b.MarkupPerRecipient,
			MarkupOnDelivery:       b.MarkupOnDelivery,
			BaseDeliveryRate:       b.BaseDeliveryRate,
			BaseDesignSpamTestRate: b.BaseDesignSpamTestRate,
		},
		ClientPays: b.ClientPays,
		Currency:   b.Currency,
	}
}
