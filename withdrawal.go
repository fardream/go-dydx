package dydx

type WithdrawResponse struct {
	Withdrawal []Withdrawal `json:"withdrawal"`
}

// Withdrawal is one type of transfer:
// https://docs.dydx.exchange/#create-withdrawal
type Withdrawal = Transfer

type FastWithdrawalParam struct {
	ClientID     string `json:"clientId"`
	ToAddress    string `json:"toAddress"`
	CreditAsset  string `json:"creditAsset"`
	CreditAmount string `json:"creditAmount"`
	DebitAmount  string `json:"debitAmount"`
	LpPositionId string `json:"lpPositionId"`
	Expiration   string `json:"expiration"`
	Signature    string `json:"signature"`
}
