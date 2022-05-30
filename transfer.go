package dydx

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	Type            string           `json:"type"`
	ID              string           `json:"id"`
	ClientID        string           `json:"clientId"`
	CreditAmount    *decimal.Decimal `json:"creditAmount,omitempty"`
	CreditAsset     string           `json:"creditAsset,omitempty"`
	DebitAmount     *decimal.Decimal `json:"debitAmount,omitempty"`
	DebitAsset      string           `json:"debitAsset,omitempty"`
	FromAddress     string           `json:"fromAddress"`
	Status          string           `json:"status"`
	ToAddress       string           `json:"toAddress,omitempty"`
	TransactionHash string           `json:"transactionHash,omitempty"`
	ConfirmedAt     *time.Time       `json:"confirmedAt,omitempty"`
	CreatedAt       time.Time        `json:"createdAt"`
}

type TransfersParam struct{}
