package krakenFuturesWSSDK

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"trade-bot/configs"
)

var (
	ErrServeWS                  = errors.New("server ws")
	ErrUnableToEstablishConnect = errors.New("unable to establish connect with kraken websocket")
	ErrUnableToWriteMessage     = errors.New("unable to write message to websocket")
	ErrSubscribeToFeed          = errors.New("subscribe to feed")
	ErrUnableToReadMessage      = errors.New("unable to read message from websocket")
	ErrCouldNotSubscribeToFeed  = errors.New("could not subscribe to feed")
	ErrConnect                  = errors.New("connect to ws")
	ErrLoopOverWS               = errors.New("loop over ws")
)

const (
	maxEstablishConnectCounter = 10
)

type WSAPI struct {
	ws             *websocket.Dialer
	wsAPIURL       string
	requestsConfig configs.KrakenWSAPIRequestsConfiguration
}

func NewWSAPI(config configs.KrakenWSConfiguration) *WSAPI {
	return &WSAPI{
		ws:             websocket.DefaultDialer,
		wsAPIURL:       config.Kraken.WSAPIURL,
		requestsConfig: config.Requests,
	}
}

// -------------------------- PUBLIC KRAKEN WEBSOCKET API ENDPOINTS -------------------------- //

func (a *WSAPI) Heartbeat(ctx context.Context) (<-chan *HeartbeatSubscriptionData, error) {
	heartbeatCh := make(chan *HeartbeatSubscriptionData)
	hearbeatArgs := KrakenSendMessageArguments{
		Event: "subscribe",
		Feed:  "heartbeat",
	}

	dataCh, errCh, err := a.serveWS(ctx, hearbeatArgs, &HeartbeatSubscriptionData{})
	if err != nil {
		return nil, err
	}

	go logErrors(errCh)
	go func() {
		defer close(heartbeatCh)
		for val := range dataCh {
			heartbeatCh <- val.(*HeartbeatSubscriptionData)
		}
	}()

	return heartbeatCh, nil
}

func (a *WSAPI) CandlesTrade(ctx context.Context, feed string, productIDs []string) (<-chan *CandlesTradeData, error) {
	candlesTradeCh := make(chan *CandlesTradeData)
	candlesArgs := KrakenSendMessageArguments{
		Event:      "subscribe",
		Feed:       feed,
		ProductIDs: productIDs,
	}

	dataCh, errCh, err := a.serveWS(ctx, candlesArgs, &CandlesTradeData{})
	if err != nil {
		return nil, err
	}

	go logErrors(errCh)
	go func() {
		defer close(candlesTradeCh)
		for val := range dataCh {
			candlesTradeCh <- val.(*CandlesTradeData)
		}
	}()

	return candlesTradeCh, nil
}

// ------------------------------------------------------------------------------------------- //

func logErrors(errCh <-chan error) {
	for val := range errCh {
		log.Warn(val)
	}
}

func (a *WSAPI) serveWS(ctx context.Context, args KrakenSendMessageArguments, typ interface{}) (<-chan interface{}, <-chan error, error) {
	conn, err := a.connect(args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", ErrServeWS, err)
	}

	dataCh, errCh := a.loopOverWS(ctx, conn, args, typ)
	return dataCh, errCh, nil
}

func (a *WSAPI) establishConnect() (*websocket.Conn, error) {
	var initResp map[string]interface{}

	for i := 0; i < maxEstablishConnectCounter; i++ {
		conn, resp, err := a.ws.Dial(a.wsAPIURL, nil)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusSwitchingProtocols {
			continue
		}

		err = conn.ReadJSON(&initResp)
		if err != nil {
			continue
		}
		if val, ok := initResp["event"]; ok && val == "info" {
			return conn, nil
		}
	}

	return nil, ErrUnableToEstablishConnect
}

func (a *WSAPI) sendEvent(conn *websocket.Conn, args KrakenSendMessageArguments) (KrakenSendMessageResponse, error) {
	if err := conn.WriteJSON(args); err != nil {
		return KrakenSendMessageResponse{}, fmt.Errorf("%s: %s: %w", ErrSubscribeToFeed, ErrUnableToWriteMessage, err)
	}

	var response KrakenSendMessageResponse
	err := conn.ReadJSON(&response)
	if err != nil {
		return response, fmt.Errorf("%s: %s: %w", ErrSubscribeToFeed, ErrUnableToReadMessage, err)
	} else if response.Event != "subscribed" {
		return response, fmt.Errorf("%s: %s", ErrSubscribeToFeed, ErrCouldNotSubscribeToFeed)
	}

	return response, nil
}

func (a *WSAPI) loopOverWS(ctx context.Context, conn *websocket.Conn, args KrakenSendMessageArguments, typ interface{}) (<-chan interface{}, <-chan error) {
	loopChan := make(chan interface{})
	errChan := make(chan error, 1)

	conn.SetReadLimit(int64(a.requestsConfig.MaxMessageSize))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(a.requestsConfig.PongWaitInSeconds)))
	})

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	go func() {
		defer close(loopChan)
		defer close(errChan)

		for {
			err := conn.ReadJSON(&typ)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					conn, err = a.connect(args)
					if err != nil {
						errChan <- fmt.Errorf("%s: %w", ErrLoopOverWS, err)
						break
					}
					continue
				}
				break
			}
			loopChan <- typ
		}
	}()

	return loopChan, errChan
}

func (a *WSAPI) connect(args KrakenSendMessageArguments) (*websocket.Conn, error) {
	conn, err := a.establishConnect()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrConnect, err)
	}

	if _, err := a.sendEvent(conn, args); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrConnect, err)
	}

	return conn, nil
}
