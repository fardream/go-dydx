package dydx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ChannelResponseHeader contains all the common information in the channel response.
type ChannelResponseHeader struct {
	// Type of the response (see ChannelResponseType*)
	Type string `json:"type"`

	// Channel
	Channel string `json:"channel"`

	// ConnectionId
	ConnectionID string `json:"connection_id,omitempty"`

	// MessageId
	MessageID int `json:"message_id,omitempty"`

	// Message contains the error message if there is an error
	Message string `json:"message,omitempty"`

	// Id is the subscribed to id, such as account, market (BTC-USD) etc.
	Id string `json:"id,omitempty"`
}

type ChannelResponse[TContents any] struct {
	ChannelResponseHeader
	Contents *TContents `json:"contents,omitempty"`
}

const (
	AccountChannel   = "v3_accounts"
	MarketsChannel   = "v3_markets"
	OrderbookChannel = "v3_orderbook"
	TradesChannel    = "v3_trades"
)

const (
	subscribeChannelRequestType   = "subscribe"
	unsubscribeChannelRequestType = "unsubscribe"
)

const (
	ChannelResponseTypeSubscribe   = "subscribed"
	ChannelResponseTypeUnsubscribe = "unsubscribed" // unused
	ChannelResponseTypeError       = "error"
	ChannelResponseTypeConnected   = "connected"
	ChannelResponseTypeChannelData = "channel_data"
)

type unsubscribeRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Id      string `json:"id,omitempty"`
}

func newUnsubscribeRequest(channel, id string) *unsubscribeRequest {
	return &unsubscribeRequest{Type: unsubscribeChannelRequestType, Channel: channel, Id: id}
}

// subscribeForType subscribes with the request and write the output to the channel.
// gorrila/websocket doesn't support context, so a separate goroutine is launched to read the data.
func subscribeForType[TData any](ctx context.Context, url string, subscribe any, unsubscribe *unsubscribeRequest, output chan<- *ChannelResponse[TData]) error {
	// wait for loop read to finish.
	var wg sync.WaitGroup
	defer wg.Wait()

	// dial and create the connection.
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 60 * time.Second,
	}

	conn, rsp, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket %s: %w\nresponse is %#v", url, err, rsp)
	}
	defer conn.Close()

	msg_chan := make(chan []byte)
	err_chan := make(chan error)

	// start looping read
	wg.Add(1)
	go func() {
		defer wg.Done()
		// close message channel
		defer close(msg_chan)
		// close error channel.
		defer close(err_chan)

		err_chan <- loopRead(ctx, conn, msg_chan)
	}()

	b, _ := json.Marshal(subscribe)
	log.Debugf("subscribe with: %s", string(b))

	if err = conn.WriteJSON(subscribe); err != nil {
		return fmt.Errorf("failed to write subscribe request %#v: %w", subscribe, err)
	}

mainloop:
	for {
		select {
		case <-ctx.Done():
			// request cancelled.
			// Close the connection.
			if err = conn.WriteJSON(unsubscribe); err != nil {
				return err
			}
			if err = conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second*2)); err != nil {
				return err
			}
			break mainloop

		case msg, ok := <-msg_chan:
			// message channel is closed.
			// break out the loop
			if !ok {
				break mainloop
			}

			// parse the response
			resp := new(ChannelResponse[TData])

			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Warnf("failed to parse data: %v", err)
				continue mainloop
			}

			if resp.Type == ChannelResponseTypeSubscribe {
				unsubscribe.Id = resp.Id
			} else if resp.Type == ChannelResponseTypeError {
				return fmt.Errorf("subscription error: %s", resp.Message)
			}

			// send the response
			select {
			case <-ctx.Done():
				break mainloop
			case output <- resp:
			}

		case err = <-err_chan:
			// the read loop errored.
			if err != nil {
				log.Warnf("err received from read loop, quit: %#v", err)
				return err
			}
		}
	}

	// drain the error channel
	return <-err_chan
}

// loopRead reads the data from the connection
func loopRead(ctx context.Context, conn *websocket.Conn, output chan<- []byte) error {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			log.Warnf("error reading websocket: %#v", err)
			return err
		}

		log.Debugf("message received from websocket: %s", string(msg))

		select {
		case output <- msg:
		case <-ctx.Done():
			return nil
		}
	}
}
