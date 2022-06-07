package dydx

import (
	"fmt"
	"strings"

	"github.com/fardream/go-dydx/heap"
)

// OrderbookProcessor maintains the state of the order book
//
// For now the bids and asks are maintained as a heap and
// it is easy to return top of the book. There doesn't look to be a
// difficult way to implement this as list as always sorted.
type OrderbookProcessor struct {
	Market string

	Bids
	Asks

	dropData bool

	Data   []*OrderbookChannelResponse
	offset int64
}

// NewOrderbookProcessor creates a orderbook processor.
// - set `dropData` to true to drop the updates.
func NewOrderbookProcessor(market string, dropData bool) *OrderbookProcessor {
	return &OrderbookProcessor{
		Market:   market,
		dropData: dropData,
		Bids:     Bids{mappedBook: mappedBook{locations: make(map[string]int)}},
		Asks:     Asks{mappedBook: mappedBook{locations: make(map[string]int)}},
	}
}

// Process a update from the orderbook
func (ob *OrderbookProcessor) Process(resp *OrderbookChannelResponse) {
	if !ob.dropData {
		ob.Data = append(ob.Data, resp)
	}

	contents := resp.Contents
	if contents == nil {
		return
	}

	if contents.Offset != nil {
		if (*contents.Offset) < ob.offset {
			return
		}
		ob.offset = *contents.Offset
	}

	ob.updateBook(contents.Bids, &ob.Bids)
	ob.updateBook(contents.Asks, &ob.Asks)
}

// updateBook updates one side of the book (bids or asks)
func (ob *OrderbookProcessor) updateBook(updates []*OrderbookOrder, book singleSideOrderbook) {
update_loop:
	for _, order := range updates {
		if order == nil {
			continue
		}
		if order.Offset != nil {
			if (*order.Offset) < ob.offset {
				continue update_loop
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
	if len(ob.Bids.Orders) > 0 {
		bid = ob.Bids.Orders[0]
	}
	var ask *OrderbookOrder
	if len(ob.Asks.Orders) > 0 {
		ask = ob.Asks.Orders[0]
	}
	return bid, ask
}

// singleSideOrderbook describe bids or asks side.
// An interface is used to remove the check of the direction in the `Less` function.
type singleSideOrderbook interface {
	// allow heap operations on this type
	heap.Interface[*OrderbookOrder]
	getPriceLevelIndex(priceStr string) (int, bool)
	updatePriceLevelSize(priceStr string, newSize Decimal)
}

var (
	_ singleSideOrderbook = (*Bids)(nil)
	_ singleSideOrderbook = (*Asks)(nil)
)

// Bids side of the book.
type Bids struct {
	mappedBook
}

func (b *Bids) Less(i, j int) bool {
	return b.Orders[i].Price.GreaterThan(b.Orders[j].Price)
}

// Asks side of the book.
type Asks struct {
	mappedBook
}

func (a *Asks) Less(i, j int) bool {
	return a.Orders[i].Price.LessThan(a.Orders[j].Price)
}

// mappedBook contains all the supporting functions for singleSideOrderbook without the less function.
type mappedBook struct {
	Orders    []*OrderbookOrder
	locations map[string]int
}

func (m *mappedBook) Len() int {
	return len(m.Orders)
}

func (m *mappedBook) Swap(i, j int) {
	m.locations[m.Orders[i].PriceString], m.locations[m.Orders[j].PriceString] = j, i
	m.Orders[i], m.Orders[j] = m.Orders[j], m.Orders[i]
}

func (m *mappedBook) Pop() *OrderbookOrder {
	if len(m.Orders) == 0 {
		return nil
	}
	order := m.Orders[len(m.Orders)-1]
	m.Orders = m.Orders[0 : len(m.Orders)-1]
	delete(m.locations, order.PriceString)

	return order
}

func (m *mappedBook) Push(order *OrderbookOrder) {
	m.Orders = append(m.Orders, order)
	m.locations[order.PriceString] = len(m.Orders) - 1
}

func (m *mappedBook) getPriceLevelIndex(priceStr string) (int, bool) {
	r, ok := m.locations[priceStr]
	return r, ok
}

func (m *mappedBook) updatePriceLevelSize(priceStr string, news_size Decimal) {
	r, ok := m.locations[priceStr]
	if ok {
		m.Orders[r].Size = news_size
	}
}

func (m *mappedBook) PrintBook() string {
	var b strings.Builder
	for index, v := range m.Orders {
		fmt.Fprintf(&b, "%d : %s @ $%s\n", index, v.Price.String(), v.Size.String())
	}
	return b.String()
}
