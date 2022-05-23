package dydx

import "time"

type TransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	Type            string    `json:"type"`
	ID              string    `json:"id"`
	ClientID        string    `json:"clientId"`
	CreditAmount    string    `json:"creditAmount"`
	CreditAsset     string    `json:"creditAsset"`
	DebitAmount     string    `json:"debitAmount"`
	DebitAsset      string    `json:"debitAsset"`
	FromAddress     string    `json:"fromAddress"`
	Status          string    `json:"status"`
	ToAddress       string    `json:"toAddress"`
	TransactionHash string    `json:"transactionHash"`
	ConfirmedAt     time.Time `json:"confirmedAt"`
	CreatedAt       time.Time `json:"createdAt"`
}

type TransfersParam struct{}
