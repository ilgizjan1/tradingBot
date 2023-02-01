package krakenFuturesSDK

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrDoRequest                = errors.New("do request")
	ErrNilTyp                   = errors.New("nil typ")
	ErrCouldNotCreateRequest    = errors.New("could not create request")
	ErrCouldNotExecuteRequest   = errors.New("could not execute request")
	ErrCouldNotReadBody         = errors.New("could not read body")
	ErrCouldNotParseContentType = errors.New("could noy parse content type")
	ErrInvalidContentType       = errors.New("invalid content type")
	ErrCouldNotUnmarshalBody    = errors.New("could not unmarshal body")
	ErrValidateSendStatus       = errors.New("validate send status")
	ErrEmptyOrderEvents         = errors.New("empty order events")
)

const (
	apiUserAgent = "Kraken GO API Agent"
)

type API struct {
	apiPublicKey  string
	apiPrivateKey string
	apiURL        string
	client        *http.Client
}

func NewAPI(apiPublicKey, apiPrivateKey string, apiURL string) *API {
	return &API{
		apiPublicKey:  apiPublicKey,
		apiPrivateKey: apiPrivateKey,
		apiURL:        apiURL,
		client:        http.DefaultClient,
	}
}

// -------------------------- PUBLIC KRAKEN API ENDPOINTS -------------------------- //

func (a *API) FeeSchedules() (*FeeSchedulesResponse, error) {
	resp, err := a.queryPublic(http.MethodGet, "/derivatives/api/v3/feeschedules", nil, &FeeSchedulesResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*FeeSchedulesResponse), nil
}

func (a *API) OrderBook(symbol string) (*OrderBookResponse, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	resp, err := a.queryPublic(http.MethodGet, "/derivatives/api/v3/orderbook", values, &OrderBookResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*OrderBookResponse), nil
}

func (a *API) Tickers() (*TickersResponse, error) {
	resp, err := a.queryPublic(http.MethodGet, "/derivatives/api/v3/tickers", nil, &TickersResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*TickersResponse), nil
}

func (a *API) Instruments() (*InstrumentsResponse, error) {
	resp, err := a.queryPublic(http.MethodGet, "/derivatives/api/v3/instruments", nil, &InstrumentsResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*InstrumentsResponse), nil
}

// --------------------------------------------------------------------------------- //

// -------------------------- PRIVATE KRAKEN API ENDPOINTS -------------------------- //

func (a *API) SendOrder(args SendOrderArguments) (*SendOrderResponse, error) {
	values := url.Values{}
	values.Add("orderType", args.OrderType)
	values.Add("symbol", args.Symbol)
	values.Add("side", args.Side)
	values.Add("size", strconv.Itoa(int(args.Size)))

	if args.LimitPrice != 0 {
		values.Add("limitPrice", strconv.FormatFloat(args.LimitPrice, 'f', 2, 64))
	}

	if args.OrderType == "stp" || args.OrderType == "take_profit" {
		if args.StopPrice != 0 {
			values.Add("stopPrice", strconv.FormatFloat(args.StopPrice, 'f', 2, 64))
		}
		if args.TriggerSignal != "" {
			values.Add("triggerSignal", args.TriggerSignal)
		}
	}

	if args.CliOrderID != "" {
		values.Add("cliOrdId", args.CliOrderID)
	}

	if args.ReduceOnly {
		values.Add("reduceOnly", "true")
	}

	resp, err := a.queryPrivate(http.MethodPost, "/derivatives/api/v3/sendorder", values, &SendOrderResponse{})
	if err != nil {
		return nil, err
	}

	if err := resp.(*SendOrderResponse).SendStatus.ValidateSendStatus(); err != nil {
		return nil, fmt.Errorf("%s: send status - %s", err, resp.(*SendOrderResponse).SendStatus.Status)
	}
	return resp.(*SendOrderResponse), nil
}

func (a *API) EditOrder(args EditOrderArguments) (*EditOrderResponse, error) {
	values := url.Values{}
	values.Add("orderId", args.OrderID)
	if args.Size != 0 {
		values.Add("size", strconv.Itoa(int(args.Size)))
	}
	if args.LimitPrice != 0 {
		values.Add("limitPrice", strconv.FormatFloat(args.LimitPrice, 'f', 2, 64))
	}
	if args.StopPrice != 0 {
		values.Add("stopPrice", strconv.FormatFloat(args.StopPrice, 'f', 2, 64))
	}
	if args.CliOrdID != "" {
		values.Add("cliOrdId", args.CliOrdID)
	}

	resp, err := a.queryPrivate(http.MethodPost, "/derivatives/api/v3/editorder", values, &EditOrderResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*EditOrderResponse), nil
}

func (a *API) CancelOrder(args CancelOrderArguments) (*CancelOrderResponse, error) {
	values := url.Values{}
	if args.OrderID != "" {
		values.Add("order_id", args.OrderID)
	}
	if args.CliOrdID != "" {
		values.Add("cliOrdId", args.CliOrdID)
	}

	resp, err := a.queryPrivate(http.MethodPost, "/derivatives/api/v3/cancelorder", values, &CancelOrderResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*CancelOrderResponse), nil
}

func (a *API) CancelAllOrders(symbol string) (*CancelAllOrdersResponse, error) {
	values := url.Values{}
	if symbol != "" {
		values.Add("symbol", symbol)
	}
	resp, err := a.queryPrivate(http.MethodPost, "/derivatives/api/v3/cancelallorders", values, &CancelAllOrdersResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*CancelAllOrdersResponse), nil
}

// ---------------------------------------------------------------------------------- //

func (s SendStatus) ValidateSendStatus() error {
	if len(s.OrderEvents) == 0 {
		return fmt.Errorf("%s: %s", ErrValidateSendStatus, ErrEmptyOrderEvents)
	}
	return nil
}

// queryPublic make request to public KrakenAPI endpoint
func (a *API) queryPublic(reqType string, endpoint string, values url.Values, typ interface{}) (interface{}, error) {
	urlPath := fmt.Sprintf("%s%s?%s", a.apiURL, endpoint, values.Encode())
	return a.doRequest(reqType, urlPath, nil, typ)
}

// queryPrivate make request to private KrakenAPI endpoint
func (a *API) queryPrivate(reqType string, endpoint string, values url.Values, typ interface{}) (interface{}, error) {
	urlPath := fmt.Sprintf("%s%s?%s", a.apiURL, endpoint, values.Encode())
	authent, err := a.createSignature(endpoint, values.Encode(), "")
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authent": authent,
		"APIKey":  a.apiPublicKey,
	}

	return a.doRequest(reqType, urlPath, headers, typ)
}

// doRequest executes HTTP Request to the KrakenAPI and returns the result
func (a *API) doRequest(reqType string, reqURL string, headers map[string]string, typ interface{}) (interface{}, error) {
	if typ == nil {
		return nil, fmt.Errorf("%s: %s", ErrDoRequest, ErrNilTyp)
	}

	// Create request
	req, err := http.NewRequest(reqType, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrDoRequest, ErrCouldNotCreateRequest, err)
	}

	req.Header.Add("User-Agent", apiUserAgent)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Execute request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrDoRequest, ErrCouldNotExecuteRequest, err)
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrDoRequest, ErrCouldNotReadBody, err)
	}

	// validate content type
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrDoRequest, ErrCouldNotParseContentType, err)
	}
	if contentType != "application/json" {
		return nil, fmt.Errorf("%s: %s: %s.\n%s",
			ErrDoRequest, ErrInvalidContentType,
			fmt.Sprintf("response content-yype is '%s', but should be 'application/json'", contentType),
			fmt.Sprintf("content: '%s'", string(body)))
	}

	err = json.Unmarshal(body, typ)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrDoRequest, ErrCouldNotUnmarshalBody, err)
	}

	return typ, nil
}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

// createSignature creates value for krakenAPI request Authent header
func (a API) createSignature(endPoint, postData, nonce string) (string, error) {
	endPoint = strings.TrimPrefix(endPoint, "/derivatives")

	message := postData + nonce + endPoint
	shaSum := getSha256([]byte(message))

	macKey, err := base64.StdEncoding.DecodeString(a.apiPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create Signature: unable to create macKey")
	}
	macSum := getHMacSha512(shaSum, macKey)
	return base64.StdEncoding.EncodeToString(macSum), nil
}
