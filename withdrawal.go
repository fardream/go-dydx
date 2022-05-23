package dydx

import "time"

type WithdrawResponse struct {
	Withdrawal []Withdrawal `json:"withdrawal"`
}

type Withdrawal struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	DebitAsset      string      `json:"debitAsset"`
	CreditAsset     string      `json:"creditAsset"`
	DebitAmount     string      `json:"debitAmount"`
	CreditAmount    string      `json:"creditAmount"`
	TransactionHash string      `json:"transactionHash"`
	Status          string      `json:"status"`
	ClientID        string      `json:"clientId"`
	FromAddress     string      `json:"fromAddress"`
	ToAddress       interface{} `json:"toAddress"`
	ConfirmedAt     interface{} `json:"confirmedAt"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type FastWithdrawParam struct {
	ClientID     string `json:"clientId"`
	ToAddress    string `json:"toAddress"`
	CreditAsset  string `json:"creditAsset"`
	CreditAmount string `json:"creditAmount"`
	DebitAmount  string `json:"debitAmount"`
	LpPositionId string `json:"lpPositionId"`
	Expiration   string `json:"expiration"`
	Signature    string `json:"signature"`
}
