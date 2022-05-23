package dydx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderBookOrder struct {
	Price  string
	Size   string
	Offset string
}

// UnmarshalJSON parse the content into an orderbook order.
// Right now the process first tries to parse the data with []string,
// if that failed, parse it with map[string]string
func (p *OrderBookOrder) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		mapper := make(map[string]string)
		if err1 := json.Unmarshal(data, &mapper); err1 != nil {
			return fmt.Errorf("failed to parse the data: as []string: %#v, and as map[string]string: %#v", err, err1)
		}
		v, ok := mapper["price"]
		if ok {
			p.Price = v
		}
		v, ok = mapper["size"]
		if ok {
			p.Size = v
		}
		v, ok = mapper["offset"]
		if ok {
			p.Offset = v
		}
	}

	l := len(s)
	switch l {
	case 2:
		p.Price = s[0]
		p.Size = s[1]
	case 3:
		p.Price = s[0]
		p.Size = s[1]
		p.Offset = s[2]
	}

	return nil
}

// OrderbookResponse is from https://docs.dydx.exchange/?json#get-orderbook
type OrderbookResponse struct {
	Offset string           `json:"offset"`
	Bids   []OrderBookOrder `json:"bids"`
	Asks   []OrderBookOrder `json:"asks"`
}

func (c *Client) GetOrderbook(ctx context.Context, market string) (*OrderbookResponse, error) {
	if market == "" {
		return nil, fmt.Errorf("market cannot be empty")
	}
	return doRequest[OrderbookResponse](ctx, c, http.MethodGet, urlJoin("orderbook", market), "", nil, true)
}
