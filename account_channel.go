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
	AccountNumber string `json:"accountNumber,omitempty"`
	ApiKey        string `json:"apiKey,omitempty"`
	Signature     string `json:"signature,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	Passphrase    string `json:"passphrase,omitempty"`
}

func newAccountChannelRequest(apiKey *ApiKey, accountNumber int) *accountChannelRequest {
	r := &accountChannelRequest{
		Type:          "subscribe",
		Channel:       AccountChannel,
		AccountNumber: fmt.Sprint(accountNumber),
		ApiKey:        apiKey.Key,
		Passphrase:    apiKey.Passphrase,
	}

	isoTimestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	r.Signature = apiKey.Sign("/ws/accounts", http.MethodGet, isoTimestamp, nil)
	r.Timestamp = isoTimestamp

	return r
}

type AccountChannelResponseContents struct {
	Orders    []Order    `json:"orders"`
	Account   Account    `json:"account"`
	Fills     []Fill     `json:"fills"`
	Accounts  []Account  `json:"accounts"`
	Transfers []Transfer `json:"transfers"`
	Positions []Position `json:"positions"`
}
type AccountChannelResponse = ChannelResponse[AccountChannelResponseContents]

func (c *Client) SubscribeAccount(ctx context.Context, accountNumber int, outputChan chan<- *AccountChannelResponse) error {
	if c.apiKey == nil {
		return fmt.Errorf("client doesn't have api key")
	}

	return subscribeForType(ctx, c.wsUrl, newAccountChannelRequest(c.apiKey, accountNumber), newUnsubscribeRequest(AccountChannel, ""), outputChan)
}
