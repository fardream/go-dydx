package dydx

// AccountInfoByMarket contains market specific information.
type AccountInfoByMarket struct {
	Market string

	// OpenPosition currently active
	OpenPosition *Position
	// All Positions
	AllPositions []*Position
	// ActiveOrders
	ActiveOrders map[string]*Order
	// ClosedOrders
	ClosedOrders map[string]*Order
	// Fills
	Fills map[string]*FillList
}

type FillList []*Fill

func NewAccountInfoByMarket(market string) *AccountInfoByMarket {
	return &AccountInfoByMarket{
		Market:       market,
		ActiveOrders: make(map[string]*Order),
		ClosedOrders: make(map[string]*Order),
		Fills:        make(map[string]*FillList),
	}
}

func (info *AccountInfoByMarket) AddPosition(position *Position) {
	if position.Status == PositionStatusOpen {
		if info.OpenPosition == position {
			return
		}
		if info.OpenPosition != nil {
			info.AllPositions = append(info.AllPositions, info.OpenPosition)
		}
		info.OpenPosition = position
	} else {
		info.AllPositions = append(info.AllPositions, position)
	}
}

func (info *AccountInfoByMarket) AddFill(fill *Fill) {
	fill_list, ok := info.Fills[fill.OrderID]
	if !ok {
		fill_list = new(FillList)
		info.Fills[fill.OrderID] = fill_list
	}
	*fill_list = append(*fill_list, fill)
}

func (info *AccountInfoByMarket) AddOrder(order *Order) {
	id := order.ID
	if order.Status == OrderStatusOpen || order.Status == OrderStatusPending {
		info.ActiveOrders[id] = order
	} else {
		delete(info.ActiveOrders, id)
		info.ClosedOrders[id] = order
	}
}

// AccountProcessor can be used to process Account Channel Updates
type AccountProcessor struct {
	// processed response.
	data []*AccountChannelResponse

	Account *Account

	Info map[string]*AccountInfoByMarket
}

func NewAccountProcessor() *AccountProcessor {
	return &AccountProcessor{
		Info: make(map[string]*AccountInfoByMarket),
	}
}

func (ap *AccountProcessor) getAccountInfoByMarket(market string) *AccountInfoByMarket {
	r, ok := ap.Info[market]
	if !ok {
		r = NewAccountInfoByMarket(market)
	}

	return r
}

// ProcessChannelResponse processes the channel responses in sequence.
func (ap *AccountProcessor) ProcessChannelResponse(resp *AccountChannelResponse) {
	// store the update
	ap.data = append(ap.data, resp)
	// contents
	contents := resp.Contents
	if contents == nil {
		return
	}

	// only update account if ap.Account is nil
	if ap.Account == nil {
		if contents.Account != nil {
			ap.Account = contents.Account
		}
	}

	for _, position := range contents.Positions {
		ap.getAccountInfoByMarket(position.Market).AddPosition(position)
	}

	for _, order := range contents.Orders {
		ap.getAccountInfoByMarket(order.Market).AddOrder(order)
	}

	for _, fill := range contents.Fills {
		ap.getAccountInfoByMarket(fill.Market).AddFill(fill)
	}
}
