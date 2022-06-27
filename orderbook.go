package dydx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OrderbookOrder is an entry on the order book, it only contains price, quantity, and potentially an offset.
type OrderbookOrder struct {
	Price  *Decimal `json:"price"`
	Size   *Decimal `json:"size"`
	Offset *int64   `json:"offset,omitempty"`

	PriceString string `json:"-"`
}

// UnmarshalJSON parse the content into an orderbook order.
// Right now the process first tries to parse the data with []string,
// if that failed, parse it with map[string]string
func (p *OrderbookOrder) UnmarshalJSON(data []byte) error {
	var s []Decimal
	if err := json.Unmarshal(data, &s); err != nil {
		mapper := make(map[string]Decimal)
		if err1 := json.Unmarshal(data, &mapper); err1 != nil {
			return fmt.Errorf("failed to parse the data: as []Decimal: %#v, and as map[string]string: %#v: %s", err, err1, string(data))
		}
		if v, ok := mapper["price"]; ok {
			p.Price = &v
		}
		if v, ok := mapper["size"]; ok {
			p.Size = &v
		}
		if v, ok := mapper["offset"]; ok {
			k, err := v.Int64()
			if err != nil {
				return fmt.Errorf("offset is not an integer: %v", v)
			}
			p.Offset = &k
		}
	} else {
		l := len(s)
		switch l {
		case 2:
			p.Price = &s[0]
			p.Size = &s[1]
		case 3:
			p.Price = &s[0]
			p.Size = &s[1]
			k, err := s[2].Int64()
			if err != nil {
				return fmt.Errorf("offset is not an integer: %v", s[2])
			}
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

func (o *OrderbookOrder) IsOtherNewerOffset(other *OrderbookOrder) bool {
	if o.Offset == nil || other.Offset == nil {
		return true
	}

	return (*o.Offset) < (*other.Offset)
}
