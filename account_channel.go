package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type accountChannelRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`

	// Private
	AccountNumber AccountNumber `json:"accountNumber"`
	ApiKey        string        `json:"apiKey,omitempty"`
	Signature     string        `json:"signature,omitempty"`
	Timestamp     string        `json:"timestamp,omitempty"`
	Passphrase    string        `json:"passphrase,omitempty"`
}

func newAccountChannelRequest(apiKey *ApiKey, accountNumber int) *accountChannelRequest {
	r := &accountChannelRequest{
		Type:          "subscribe",
		Channel:       AccountChannel,
		AccountNumber: AccountNumber(accountNumber),
		ApiKey:        apiKey.Key,
		Passphrase:    apiKey.Passphrase,
	}

	isoTimestamp := GetIsoDateStr(time.Now())
	r.Signature = apiKey.Sign("/ws/accounts", http.MethodGet, isoTimestamp, nil)
	r.Timestamp = isoTimestamp

	return r
}

type AccountChannelResponseContents struct {
	Account   *Account    `json:"account,omitempty"`
	Orders    []*Order    `json:"orders,omitempty"`
	Fills     []*Fill     `json:"fills,omitempty"`
	Accounts  []*Account  `json:"accounts,omitempty"`
	Transfers []*Transfer `json:"transfers,omitempty"`
	Positions []*Position `json:"positions,omitempty"`
}

type AccountChannelResponse = ChannelResponse[AccountChannelResponseContents]

// SubscribeAccount gets the accounts update
// It will feed the account update in sequence into the channnel provided. It returns after the subscription is done and closed.
func (c *Client) SubscribeAccount(ctx context.Context, accountNumber int, outputChan chan<- *AccountChannelResponse) error {
	if c.apiKey == nil {
		return fmt.Errorf("client doesn't have api key")
	}

	return subscribeForType(ctx, c.wsUrl, newAccountChannelRequest(c.apiKey, accountNumber), newUnsubscribeRequest(AccountChannel, ""), outputChan)
}
