package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"trade-bot/configs"
)

var (
	ErrDoWS                       = errors.New("do ws")
	ErrUnableToConnectToWebsocket = errors.New("unable to connect to websocket")
)

const maxEstablishConnectCounter = 5

type ClientActions interface {
	NewRequest(method, path, jwtToken string, body interface{}) (*http.Request, error)
	NewWsRequest(path, jwtToken string) (*http.Request, error)
	Do(req *http.Request, typ interface{}) (*http.Response, error)
	DoWS(req *http.Request, typ interface{}) (*websocket.Conn, error)
	LoopOverWS(conn *websocket.Conn, typ interface{}) (<-chan interface{}, <-chan error)
}

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
	ws         *websocket.Dialer
}

func NewClient(config configs.ClientConfiguration) (*Client, error) {
	urlValue, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		BaseURL:    urlValue,
		UserAgent:  "go client user agent",
		httpClient: http.DefaultClient,
		ws:         websocket.DefaultDialer,
	}

	return c, nil
}

func (c *Client) NewWsRequest(path, jwtToken string) (*http.Request, error) {
	req, err := c.NewRequest(http.MethodGet, path, jwtToken, nil)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = "ws"
	return req, nil
}

func (c *Client) DoWS(req *http.Request, typ interface{}) (*websocket.Conn, error) {
	for i := 0; i < maxEstablishConnectCounter; i++ {
		conn, resp, err := c.ws.Dial(req.URL.String(), req.Header)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusSwitchingProtocols {
			continue
		}

		if err := conn.WriteJSON(typ); err != nil {
			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("%s: %s", ErrDoWS, ErrUnableToConnectToWebsocket)
}

func (c *Client) LoopOverWS(conn *websocket.Conn, typ interface{}) (<-chan interface{}, <-chan error) {
	loopChan := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(loopChan)
		defer close(errCh)

		for {
			if err := conn.ReadJSON(&typ); err != nil {
				return
			}
			loopChan <- typ
		}
	}()

	return loopChan, errCh
}

func (c *Client) NewRequest(method, path, jwtToken string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if jwtToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, typ interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(typ)
	return resp, err
}
