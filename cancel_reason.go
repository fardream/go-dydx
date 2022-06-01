package dydx

import "encoding/json"

// CancelReason is the reason for order cancellation.
// See here https://docs.dydx.exchange/#get-orders (section "Cancel Reasons")
type CancelReason string

const (
	CancelReasonUndercollateralized CancelReason = "UNDERCOLLATERALIZED"   // Order would have led to an undercollateralized state for the user.
	CancelReasonExpired             CancelReason = "EXPIRED"               // Order expired.
	CancelReasonUserCancelled       CancelReason = "USER_CANCELED"         // Order was canceled by the user.
	CancelReasonSelfTrade           CancelReason = "SELF_TRADE"            // Order would have resulted in a self trade for the user.
	CancelReasonFailed              CancelReason = "FAILED"                // An internal issue caused the order to be canceled.
	CancelReasonCouldNotFill        CancelReason = "COULD_NOT_FILL"        // A FOK or IOC order could not be fully filled.
	CancelReasonPostOnlyWouldCross  CancelReason = "POST_ONLY_WOULD_CROSS" // A post-only order would cross the orderbook.
)

var cancelReasons []string = []string{
	string(CancelReasonUndercollateralized),
	string(CancelReasonExpired),
	string(CancelReasonUserCancelled),
	string(CancelReasonSelfTrade),
	string(CancelReasonFailed),
	string(CancelReasonCouldNotFill),
	string(CancelReasonPostOnlyWouldCross),
}

func GetCancelReason(input string) (CancelReason, error) {
	return getProperStringEnum[CancelReason](input, cancelReasons, "CancelReason")
}

var (
	_ json.Marshaler   = (*CancelReason)(nil)
	_ json.Unmarshaler = (*CancelReason)(nil)
)

func (ot *CancelReason) UnmarshalJSON(input []byte) error {
	return unmarshalJsonForStringEnum(input, "CancelReason", ot, cancelReasons)
}

func (ot CancelReason) MarshalJSON() ([]byte, error) {
	return marshalJsonForStringEnum(ot, "CancelReason", cancelReasons)
}
