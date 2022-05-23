package dydx

import (
	"context"
	"net/http"
	"time"
)

type UsersResponse struct {
	User User `json:"user"`
}

type User struct {
	PublicID        string `json:"publicId"`
	EthereumAddress string `json:"ethereumAddress"`
	IsRegistered    bool   `json:"isRegistered"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	UserData        struct {
		WalletType  string `json:"walletType"`
		Preferences struct {
			SaveOrderAmount  bool `json:"saveOrderAmount"`
			UserTradeOptions struct {
				Limit struct {
					PostOnlyChecked           bool   `json:"postOnlyChecked"`
					GoodTilTimeInput          string `json:"goodTilTimeInput"`
					GoodTilTimeTimescale      string `json:"goodTilTimeTimescale"`
					SelectedTimeInForceOption string `json:"selectedTimeInForceOption"`
				} `json:"LIMIT"`
				Market struct {
					PostOnlyChecked           bool   `json:"postOnlyChecked"`
					GoodTilTimeInput          string `json:"goodTilTimeInput"`
					GoodTilTimeTimescale      string `json:"goodTilTimeTimescale"`
					SelectedTimeInForceOption string `json:"selectedTimeInForceOption"`
				} `json:"MARKET"`
				StopLimit struct {
					PostOnlyChecked           bool   `json:"postOnlyChecked"`
					GoodTilTimeInput          string `json:"goodTilTimeInput"`
					GoodTilTimeTimescale      string `json:"goodTilTimeTimescale"`
					SelectedTimeInForceOption string `json:"selectedTimeInForceOption"`
				} `json:"STOP_LIMIT"`
				TakeProfit struct {
					PostOnlyChecked           bool   `json:"postOnlyChecked"`
					GoodTilTimeInput          string `json:"goodTilTimeInput"`
					GoodTilTimeTimescale      string `json:"goodTilTimeTimescale"`
					SelectedTimeInForceOption string `json:"selectedTimeInForceOption"`
				} `json:"TAKE_PROFIT"`
				LastPlacedTradeType string `json:"lastPlacedTradeType"`
			} `json:"userTradeOptions"`
			PopUpNotifications      bool      `json:"popUpNotifications"`
			OrderbookAnimations     bool      `json:"orderbookAnimations"`
			OneTimeNotifications    []string  `json:"oneTimeNotifications"`
			LeaguesCurrentStartDate time.Time `json:"leaguesCurrentStartDate"`
		} `json:"preferences"`
		Notifications struct {
			Trade struct {
				Email bool `json:"email"`
			} `json:"trade"`
			Deposit struct {
				Email bool `json:"email"`
			} `json:"deposit"`
			Transfer struct {
				Email bool `json:"email"`
			} `json:"transfer"`
			Marketing struct {
				Email bool `json:"email"`
			} `json:"marketing"`
			Withdrawal struct {
				Email bool `json:"email"`
			} `json:"withdrawal"`
			Liquidation struct {
				Email bool `json:"email"`
			} `json:"liquidation"`
			FundingPayment struct {
				Email bool `json:"email"`
			} `json:"funding_payment"`
		} `json:"notifications"`
		StarredMarkets []interface{} `json:"starredMarkets"`
	} `json:"userData"`
	MakerFeeRate                 string `json:"makerFeeRate"`
	TakerFeeRate                 string `json:"takerFeeRate"`
	MakerVolume30D               string `json:"makerVolume30D"`
	TakerVolume30D               string `json:"takerVolume30D"`
	Fees30D                      string `json:"fees30D"`
	ReferredByAffiliateLink      string `json:"referredByAffiliateLink"`
	IsSharingUsername            bool   `json:"isSharingUsername"`
	IsSharingAddress             bool   `json:"isSharingAddress"`
	DydxTokenBalance             string `json:"dydxTokenBalance"`
	StakedDydxTokenBalance       string `json:"stakedDydxTokenBalance"`
	ActiveStakedDydxTokenBalance string `json:"activeStakedDydxTokenBalance"`
	IsEmailVerified              bool   `json:"isEmailVerified"`
	Country                      any    `json:"country"`
	HedgiesHeld                  []any  `json:"hedgiesHeld"`
}

func (c *Client) GetUser(ctx context.Context) (*UsersResponse, error) {
	return doRequest[UsersResponse](ctx, c, http.MethodGet, "users", "", nil, false)
}
