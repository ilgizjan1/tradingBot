package krakenFuturesSDK

const SellSide = "sell"
const BuySide = "buy"

type SendOrderStatus string

func (s SendOrderStatus) IsSuccessStatus() bool {
	statuses := map[string]struct{}{"placed": {}, "cancelled": {}}
	if _, ok := statuses[string(s)]; ok {
		return true
	}
	return false
}

type EditOrderStatus string

func (s EditOrderStatus) IsSuccessStatus() bool {
	statuses := map[string]struct{}{"edited": {}}
	if _, ok := statuses[string(s)]; ok {
		return true
	}
	return false
}

type CancelOrderStatus string

func (s CancelOrderStatus) IsSuccessStatus() bool {
	statuses := map[string]struct{}{"cancelled": {}}
	if _, ok := statuses[string(s)]; ok {
		return true
	}
	return false
}

type CancelAllOrdersStatus string

func (s CancelAllOrdersStatus) IsSuccessStatus() bool {
	statuses := map[string]struct{}{"cancelled": {}}
	if _, ok := statuses[string(s)]; ok {
		return true
	}
	return false
}

// KrakenErrorResponse wraps the Kraken API JSON error response
type KrakenErrorResponse struct {
	Result     string `json:"result,omitempty"`
	ServerTime string `json:"serverTime,omitempty"`
	Error      string `json:"error,omitempty"`
}

// -------------------------- PUBLIC KRAKEN API ENDPOINTS DATA -------------------------- //

// FeeSchedulesResponse wraps the Kraken API JSON FeeSchedules method
type FeeSchedulesResponse struct {
	KrakenErrorResponse
	FeeSchedules []FeeSchedules `json:"feeSchedules,omitempty"`
}

// OrderBookResponse wraps the Kraken API JSON OrderBook method
type OrderBookResponse struct {
	KrakenErrorResponse
	OrderBook OrderBook `json:"orderBook,omitempty"`
}

// TickersResponse wraps the Kraken API JSON Tickers method
type TickersResponse struct {
	KrakenErrorResponse
	Tickers []Ticker `json:"tickers,omitempty"`
}

// InstrumentsResponse wraps the Kraken API JSON Instruments method
type InstrumentsResponse struct {
	KrakenErrorResponse
	Instruments []Instrument `json:"instruments,omitempty"`
}

// --------------------------------------------------------------------------------------- //

// -------------------------- PRIVATE KRAKEN API ENDPOINTS DATA -------------------------- //

type SendOrderResponse struct {
	KrakenErrorResponse
	SendStatus SendStatus `json:"sendStatus,omitempty"`
}

type SendOrderArguments struct {
	OrderType     string  `json:"order_type" binding:"required"`
	Symbol        string  `json:"symbol" binding:"required"`
	Side          string  `json:"side" binding:"required"`
	Size          uint    `json:"size" binding:"required"`
	LimitPrice    float64 `json:"limit_price"`
	StopPrice     float64 `json:"stop_price"`
	TriggerSignal string  `json:"trigger_signal"`
	CliOrderID    string  `json:"cli_order_id"`
	ReduceOnly    bool    `json:"reduce_only"`
}

func (s *SendOrderArguments) ChangeToOpositeOrderSide() {
	if s.Side == BuySide {
		s.Side = SellSide
	} else {
		s.Side = BuySide
	}
}

type EditOrderResponse struct {
	KrakenErrorResponse
	EditStatus EditStatus `json:"editStatus,omitempty"`
}

type EditOrderArguments struct {
	OrderID    string
	Size       uint
	LimitPrice float64
	StopPrice  float64
	CliOrdID   string
}

type CancelOrderResponse struct {
	KrakenErrorResponse
	CancelStatus CancelStatus `json:"cancelStatus,omitempty"`
}

type CancelOrderArguments struct {
	OrderID  string
	CliOrdID string
}

type CancelAllOrdersResponse struct {
	KrakenErrorResponse
	CancelStatus CancelAllStatus `json:"cancelStatus,omitempty"`
}

// --------------------------------------------------------------------------------------- //

type CancelStatus struct {
	Status       CancelOrderStatus `json:"status"`
	OrderID      string            `json:"order_id"`
	CliOrdID     string            `json:"cliOrdId,omitempty"`
	ReceivedTime string            `json:"receivedTime"`
	OrderEvents  []OrderEvent      `json:"orderEvents,omitempty"`
}

type EditStatus struct {
	OrderID      string          `json:"orderId"`
	CliOrderID   string          `json:"cliOrderId,omitempty"`
	ReceivedTime string          `json:"receivedTime"`
	Status       EditOrderStatus `json:"status"`
	OrderEvents  []OrderEvent    `json:"orderEvents,omitempty"`
}

type SendStatus struct {
	OrderID      string          `json:"order_id,omitempty"`
	CliOrderID   string          `json:"cliOrderId,omitempty"`
	Status       SendOrderStatus `json:"status,omitempty"`
	ReceivedTime string          `json:"receivedTime,omitempty"`
	OrderEvents  []OrderEvent    `json:"orderEvents,omitempty"`
}

type CancelAllStatus struct {
	ReceivedTime    string                `json:"receivedTime,omitempty"`
	CancelOnly      string                `json:"cancelOnly,omitempty"`
	Status          CancelAllOrdersStatus `json:"status,omitempty"`
	CancelledOrders []CanceledOrder       `json:"cancelledOrders,omitempty"`
	OrderEvents     []OrderEvent          `json:"orderEvents,omitempty"`
}

type CanceledOrder struct {
	OrderID  string `json:"order_id,omitempty"`
	CliOrdID string `json:"cliOrdId,omitempty"`
}

type OrderEvent struct {
	Type                string  `json:"type,omitempty"`
	ReducedQuantity     int     `json:"reducedQuantity,omitempty"`
	Order               Order   `json:"order,omitempty"`
	UID                 string  `json:"uid,omitempty"`
	Old                 Order   `json:"old,omitempty"`
	New                 Order   `json:"new,omitempty"`
	Reason              string  `json:"reason,omitempty"`
	Amount              int     `json:"amount,omitempty"`
	Price               float64 `json:"price,omitempty"`
	ExecutionID         string  `json:"executionId,omitempty"`
	TakeReducedQuantity int     `json:"takeReducedQuantity,omitempty"`
	OrderPriorEdit      Order   `json:"orderPriorEdit,omitempty"`
	OrderPriorExecution Order   `json:"orderPriorExecution,omitempty"`
}

type Order struct {
	OrderID             string  `json:"orderId,omitempty"`
	CliOrderID          string  `json:"cliOrdID,omitempty"`
	ReduceOnly          bool    `json:"reduceOnly"`
	Symbol              string  `json:"symbol,omitempty"`
	Quantity            float64 `json:"quantity,omitempty"`
	Side                string  `json:"side,omitempty"`
	LimitPrice          float64 `json:"limitPrice,omitempty"`
	StopPrice           float64 `json:"stopPrice,omitempty"`
	Filled              float64 `json:"filled"`
	Type                string  `json:"type,omitempty"`
	Timestamp           string  `json:"timestamp,omitempty"`
	LastUpdateTimestamp string  `json:"lastUpdateTimestamp,omitempty"`
}

type Instrument struct {
	Symbol          string        `json:"symbol"`
	Type            string        `json:"type"`
	Tradeable       bool          `json:"tradeable"`
	Underlying      string        `json:"underlying,omitempty"`
	LastTradingTime string        `json:"lastTradingTime,omitempty"`
	TickSize        float64       `json:"tickSize,omitempty"`
	ContractSize    int           `json:"contractSize,omitempty"`
	MarginLevels    []MarginLevel `json:"marginLevels,omitempty"`
}

type MarginLevel struct {
	Contracts         int     `json:"contracts"`
	InitialMargin     float64 `json:"initialMargin"`
	MaintenanceMargin float64 `json:"maintenanceMargin"`
}

type OrderBook struct {
	Bids [][2]float64 `json:"bids"`
	Asks [][2]float64 `json:"asks"`
}

type FeeSchedules struct {
	Name  string `json:"name"`
	UID   string `json:"uid"`
	Tiers []Tier `json:"tiers"`
}

type Tier struct {
	MakerFee  float64 `json:"makerFee"`
	TakerFee  float64 `json:"takerFee"`
	UsdVolume float64 `json:"usdVolume"`
}

type Ticker struct {
	Tag                   string  `json:"tag,omitempty"`
	Pair                  string  `json:"pair,omitempty"`
	Symbol                string  `json:"symbol,omitempty"`
	MarkPrice             float64 `json:"markPrice,omitempty"`
	Bid                   float64 `json:"bid,omitempty"`
	BidSize               int     `json:"bidSize,omitempty"`
	Ask                   float64 `json:"ask,omitempty"`
	AskSize               int     `json:"askSize,omitempty"`
	Vol24h                int     `json:"vol24h,omitempty"`
	OpenInterest          float64 `json:"openInterest,omitempty"`
	Open24H               float64 `json:"open24h,omitempty"`
	Last                  float64 `json:"last,omitempty"`
	LastTime              string  `json:"lastTime,omitempty"`
	LastSize              int     `json:"lastSize,omitempty"`
	Suspended             bool    `json:"suspended,omitempty"`
	FundingRate           float64 `json:"funding_rate,omitempty"`
	FundingRatePrediction float64 `json:"funding_rate_prediction,omitempty"`
}
