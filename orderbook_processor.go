package dydx

import (
	"github.com/shopspring/decimal"

	"github.com/fardream/go-dydx/heap"
)

// OrderbookProcessor maintains the state of the order book
//
// For now the bids and asks are maintained as a heap and
// it is easy to return top of the book. There doesn't look to be a
// difficult way to implement this as list as always sorted.
type OrderbookProcessor struct {
	Market string

	bids
	asks

	dropData bool

	data   []*OrderbookChannelResponse
	offset int64
}

// NewOrderbookProcessor creates a bookorder processor.
// - set `dropData` to true to drop the updates.
func NewOrderbookProcessor(market string, dropData bool) *OrderbookProcessor {
	return &OrderbookProcessor{
		Market:   market,
		dropData: dropData,
		bids:     bids{mappedBook: mappedBook{locations: make(map[string]int)}},
		asks:     asks{mappedBook: mappedBook{locations: make(map[string]int)}},
	}
}

// Process a update from the orderbook
func (ob *OrderbookProcessor) Process(resp *OrderbookChannelResponse) {
	if !ob.dropData {
		ob.data = append(ob.data, resp)
	}

	contents := resp.Contents
	if contents == nil {
		return
	}

	if contents.Offset != nil {
		if (*contents.Offset) <= ob.offset {
			return
		}
		ob.offset = *contents.Offset
	}

	ob.updateBook(contents.Bids, &ob.bids)
	ob.updateBook(contents.Asks, &ob.asks)
}

// updateBook updates one side of the book (bids or asks)
func (ob *OrderbookProcessor) updateBook(updates []*OrderbookOrder, book singleSideOrderbook) {
updateloop:
	for _, order := range updates {
		if order == nil {
			continue
		}
		if order.Offset != nil {
			if (*order.Offset) < ob.offset {
				continue updateloop
			}
			ob.offset = *order.Offset
		}
		updatePriceLevel(book, order)
	}
}

// updatePriceLevel update one price level
func updatePriceLevel[T singleSideOrderbook](ob T, order *OrderbookOrder) {
	index, ok := ob.getPriceLevelIndex(order.PriceString)
	switch {
	case order.Size.IsZero() && ok:
		heap.Remove[T, *OrderbookOrder](ob, index)
	case !order.Size.IsZero() && ok:
		ob.updatePriceLevelSize(order.PriceString, order.Size)
	case !order.Size.IsZero() && !ok:
		heap.Push(ob, order)
	}
}

// BookTop returns the best bid and ask of the book. nil if the side of the book is empty.
func (ob *OrderbookProcessor) BookTop() (*OrderbookOrder, *OrderbookOrder) {
	var bid *OrderbookOrder
	if len(ob.bids.orders) > 0 {
		bid = ob.bids.orders[0]
	}
	var ask *OrderbookOrder
	if len(ob.asks.orders) > 0 {
		ask = ob.asks.orders[0]
	}
	return bid, ask
}

// singleSideOrderbook describe bids or asks side.
// An interface is used to remove the check of the direction in the `Less` function.
type singleSideOrderbook interface {
	// allow heap operations on this type
	heap.Interface[*OrderbookOrder]
	getPriceLevelIndex(pricestr string) (int, bool)
	updatePriceLevelSize(pricestr string, news_size decimal.Decimal)
}

var (
	_ singleSideOrderbook = (*bids)(nil)
	_ singleSideOrderbook = (*asks)(nil)
)

type bids struct {
	mappedBook
}

func (b *bids) Less(i, j int) bool {
	return b.orders[i].Price.GreaterThan(b.orders[j].Price)
}

type asks struct {
	mappedBook
}

func (a *asks) Less(i, j int) bool {
	return a.orders[i].Price.LessThan(a.orders[j].Price)
}

// mappedBook contains all the supporting functions for singleSideOrderbook without the less function.
type mappedBook struct {
	orders    []*OrderbookOrder
	locations map[string]int
}

func (m *mappedBook) Len() int {
	return len(m.orders)
}

func (m *mappedBook) Swap(i, j int) {
	m.locations[m.orders[i].PriceString], m.locations[m.orders[j].PriceString] = j, i
	m.orders[i], m.orders[j] = m.orders[j], m.orders[i]
}

func (m *mappedBook) Pop() *OrderbookOrder {
	if len(m.orders) == 0 {
		return nil
	}
	order := m.orders[len(m.orders)-1]
	delete(m.locations, order.PriceString)

	return order
}

func (m *mappedBook) Push(order *OrderbookOrder) {
	m.orders = append(m.orders, order)
	m.locations[order.PriceString] = len(m.orders) - 1
}

func (m *mappedBook) getPriceLevelIndex(pricestr string) (int, bool) {
	r, ok := m.locations[pricestr]
	return r, ok
}

func (m *mappedBook) updatePriceLevelSize(pricestr string, news_size decimal.Decimal) {
	r, ok := m.locations[pricestr]
	if ok {
		m.orders[r].Size = news_size
	}
}
