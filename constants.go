package dydx

const (
	ApiHostMainnet = "https://api.dydx.exchange"
	ApiHostRopsten = "https://api.stage.dydx.exchange"
	WsHostMainnet  = "wss://api.dydx.exchange/v3/ws"
	WsHostRopsten  = "wss://api.stage.dydx.exchange/v3/ws"
)

const (
	SignatureTypeNoPrepend   = 0
	SignatureTypeDecimal     = 1
	SignatureTypeHexadecimal = 2
)

const (
	Domain                       = "dYdX"
	Version                      = "1.0"
	Eip712DomainStringNoContract = "EIP712Domain(string name,string version,uint256 chainId)"
)

const (
	NetworkIdMainnet = 1
	NetworkIdRopsten = 3
)

const (
	OPEN        = "OPEN"
	CLOSED      = "CLOSED"
	LIQUIDATED  = "LIQUIDATED"
	LIQUIDATION = "LIQUIDATION"
)

const (
	OrderTypeMarket       = "MARKET"
	OrderTypeLimit        = "LIMIT"
	OrderTypeStop         = "STOP"
	OrderTypeTrailingStop = "TRAILING_STOP"
	OrderTypeTakingProfit = "TAKE_PROFIT"
)

const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

const (
	TimeInForceGtt = "GTT" // Good til time
	TimeInForceFok = "FOK" // Fill or Kill
	TimeInForceIoc = "IOC" // Immediate or Cancel
)

const (
	Resolution1D     = "1DAY"
	Resolution4HOURS = "4HOURS"
	Resolution1HOUR  = "1HOUR"
	Resolution30MINS = "30MINS"
	Resolution15MINS = "15MINS"
	Resolution5MINS  = "5MINS"
	Resolution1MIN   = "1MIN"
)
