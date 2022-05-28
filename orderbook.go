package dydx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type OrderbookOrder struct {
	Price  decimal.Decimal `json:"price"`
	Size   decimal.Decimal `json:"size"`
	Offset *int64          `json:"offset,omitempty"`

	PriceString string `json:"-"`
}

// UnmarshalJSON parse the content into an orderbook order.
// Right now the process first tries to parse the data with []string,
// if that failed, parse it with map[string]string
func (p *OrderbookOrder) UnmarshalJSON(data []byte) error {
	var s []decimal.Decimal
	if err := json.Unmarshal(data, &s); err != nil {
		mapper := make(map[string]decimal.Decimal)
		if err1 := json.Unmarshal(data, &mapper); err1 != nil {
			return fmt.Errorf("failed to parse the data: as []string: %#v, and as map[string]string: %#v", err, err1)
		}
		if v, ok := mapper["price"]; ok {
			p.Price = v
		}
		if v, ok := mapper["size"]; ok {
			p.Size = v
		}
		if v, ok := mapper["offset"]; ok {
			k := v.IntPart()
			p.Offset = &k
		}
	} else {
		l := len(s)
		switch l {
		case 2:
			p.Price = s[0]
			p.Size = s[1]
		case 3:
			p.Price = s[0]
			p.Size = s[1]
			k := s[2].IntPart()
			p.Offset = &k
		}
	}

	p.PriceString = p.Price.String()

	return nil
}

// OrderbookResponse is from https://docs.dydx.exchange/?json#get-orderbook
type OrderbookResponse struct {
	Offset *int64            `json:"offset,string,omitempty"`
	Bids   []*OrderbookOrder `json:"bids"`
	Asks   []*OrderbookOrder `json:"asks"`
}

func (c *Client) GetOrderbook(ctx context.Context, market string) (*OrderbookResponse, error) {
	if market == "" {
		return nil, fmt.Errorf("market cannot be empty")
	}
	return doRequest[OrderbookResponse](ctx, c, http.MethodGet, urlJoin("orderbook", market), "", nil, true)
}
