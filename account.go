package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type AccountResponse struct {
	Account Account `json:"account"`
}

type Account struct {
	PositionId         int64               `json:"positionId,string"`
	ID                 string              `json:"id"`
	StarkKey           string              `json:"starkKey"`
	Equity             Decimal             `json:"equity"`
	FreeCollateral     Decimal             `json:"freeCollateral"`
	QuoteBalance       Decimal             `json:"quoteBalance"`
	PendingDeposits    string              `json:"pendingDeposits"`
	PendingWithdrawals string              `json:"pendingWithdrawals"`
	AccountNumber      JsonInt             `json:"accountNumber"`
	OpenPositions      map[string]Position `json:"openPositions,omitempty"`
	CreatedAt          time.Time           `json:"createdAt"`
}

type AccountsResponse struct {
	Accounts []*Account `json:"accounts"`
}

// GetAccounts implements https://docs.dydx.exchange/#get-accounts, it gets all accounts
func (c *Client) GetAccounts(ctx context.Context) (*AccountsResponse, error) {
	return doRequest[AccountsResponse](ctx, c, http.MethodGet, "accounts", "", nil, false)
}

// GetAccount with a specific id, implements https://docs.dydx.exchange/#get-account
func (c *Client) GetAccount(ctx context.Context, id string) (*AccountResponse, error) {
	path := fmt.Sprintf("accounts/%s", id)
	return doRequest[AccountResponse](ctx, c, http.MethodGet, path, "", nil, false)
}
