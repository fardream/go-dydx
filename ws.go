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

type ChannelResponseHeader struct {
	Type         string `json:"type"`
	Channel      string `json:"channel"`
	ConnectionID string `json:"connection_id,omitempty"`
	MessageID    int    `json:"message_id,omitempty"`
	Id           string `json:"id,omitempty"`
}

type ChannelResponse[TContents any] struct {
	ChannelResponseHeader
	Contents *TContents `json:"contents,omitempty"`
}

const (
	UndefinedChannel = "undefined"
	ErrorChannel     = "error"
	AccountChannel   = "v3_accounts"
	MarketsChannel   = "v3_markets"
	OrderbookChannel = "v3_orderbook"
	TradesChannel    = "v3_trades"
)

const (
	subscribeChannelRequestType   = "subscribe"
	unsubscribeChannelRequestType = "unsubscribe"
	subscribeConfirmationType     = "subscribed"
	unsubscribeConfirmationType   = "unsubscribed" // unused
)

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

		log.Debugf("subscribe: msg: %s", string(msg))

		select {
		case output <- msg:
		case <-ctx.Done():
			return nil
		}
	}
}

type unsubscribeRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Id      string `json:"id,omitempty"`
}

func newUnsubscribeRequest(channel, id string) *unsubscribeRequest {
	return &unsubscribeRequest{Type: unsubscribeChannelRequestType, Channel: channel, Id: id}
}

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
		defer close(msg_chan)
		defer close(err_chan)

		err_chan <- loopRead(ctx, conn, msg_chan)
	}()
	if err = conn.WriteJSON(subscribe); err != nil {
		return fmt.Errorf("failed to write subscribe request %#v: %w", subscribe, err)
	}

mainloop:
	for {
		select {
		case <-ctx.Done():
			break mainloop
		case msg, ok := <-msg_chan:
			if !ok {
				break mainloop
			}

			resp := new(ChannelResponse[TData])
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Warnf("failed to parse data: %v", err)
				continue mainloop
			}
			if resp.Type == subscribeConfirmationType {
				unsubscribe.Id = resp.Id
			}
			select {
			case <-ctx.Done():
				break mainloop
			case output <- resp:
			}
		}
	}

	if err = conn.WriteJSON(unsubscribe); err != nil {
		return err
	}

	if err = conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second*2)); err != nil {
		return err
	}

	return <-err_chan
}
