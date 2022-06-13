package dydx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

// CreateUserParam contains the parameters to create a new user.
// Note ethereumAddress is actually not part of the api according to python implementation:
// https://github.com/dydxprotocol/dydx-v3-python/blob/914fc66e542d82080702e03f6ad078ca2901bb46/dydx3/modules/onboarding.py#L103-L111
type CreateUserParam struct {
	StarkPublicKey            string `json:"starkKey"`
	StarkPublicKeyYCoordinate string `json:"starkKeyYCoordinate"`
	// ethereumAddress, even though listed on the documentation, is not part of the request.
	EthereumAddress         string `json:"-"`
	ReferredByAffiliateLink string `json:"referredByAffiliateLink,omitempty"`
	Country                 string `json:"country,omitempty"`
}

type CreateUserResponse struct {
	User    *User    `json:"user,omitempty"`
	ApiKey  *ApiKey  `json:"apiKey,omitempty"`
	Account *Account `json:"account,omitempty"`
}

func (c *Client) checkCreateUserParam(param *CreateUserParam) (*CreateUserParam, error) {
	p := &CreateUserParam{}
	if param != nil {
		*p = *param
	}

	if len(p.EthereumAddress) == 0 {
		p.EthereumAddress = c.ethAddress
	}
	if len(p.StarkPublicKey) == 0 {
		if c.starkKey == nil {
			return nil, fmt.Errorf("parameter doesn't have stark public key and client doesn't have it either")
		}
		p.StarkPublicKey = c.starkKey.PublicKey
	}
	if len(p.StarkPublicKeyYCoordinate) == 0 {
		if c.starkKey == nil {
			return nil, fmt.Errorf("parameter doesn't have stark public key y coordinate and client doesn't have it either")
		}
		p.StarkPublicKeyYCoordinate = c.starkKey.PublicKeyYCoordinate
	}

	if len(p.EthereumAddress) == 0 {
		return nil, fmt.Errorf("%#v doesn't have ethereum address", *p)
	}
	if len(p.StarkPublicKey) == 0 {
		return nil, fmt.Errorf("%#v doesn't have stark public key", *p)
	}
	if len(p.StarkPublicKeyYCoordinate) == 0 {
		return nil, fmt.Errorf("%#v doesn't have stark public key y coordinate", *p)
	}

	return p, nil
}

// CreateUser creates a new user on dydx.
// See here: https://docs.dydx.exchange/#onboarding
// Note this is re-produced from python version
// https://github.com/dydxprotocol/dydx-v3-python/blob/914fc66e542d82080702e03f6ad078ca2901bb46/dydx3/modules/onboarding.py#L32-L112
func (c *Client) CreateUser(ctx context.Context, signer SignTypedData, param *CreateUserParam) (*CreateUserResponse, error) {
	p, err := c.checkCreateUserParam(param)
	if err != nil {
		return nil, fmt.Errorf("parameter error: %w", err)
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request to json: %w", err)
	}

	onboardingTypedData := getOnboardingTypedData(c.networkId == NetworkIdMainnet, onboardingAction)
	signature, err := signer.EthSignTypedData(onboardingTypedData)
	if err != nil {
		return nil, fmt.Errorf("failed to sign typed data %#v: %w", onboardingTypedData, err)
	}
	signature = append(signature, 0)

	full_path := urlJoin(c.rpcUrl, "v3/onboarding")

	log.Debugf("sending %s request to %s", http.MethodPost, full_path)

	timeout_ctx, cancel := context.WithTimeout(ctx, c.timeOut)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout_ctx, http.MethodPost, full_path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("DYDX-SIGNATURE", hexutil.Encode(signature))
	req.Header.Add("DYDX-ETHEREUM-ADDRESS", p.EthereumAddress)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &DydxError{HttpStatusCode: resp.StatusCode, Message: resp.Status, Body: msg}
	}

	log.Debugf("response from remote: %s", msg)

	r := new(CreateUserResponse)

	if err := json.Unmarshal(msg, r); err != nil {
		log.Warnf("failed to unmarshal body:\n%s", msg)
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return r, nil
}
