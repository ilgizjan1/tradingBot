package krakenFuturesWSSDK

const OneMinuteCandlesFeed = "candles_trade_1m"

// -------------------------- PUBLIC KRAKEN WEBSOCKET API DATA -------------------------- //

type KrakenSendMessageArguments struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIDs []string `json:"product_ids"`
}

type KrakenSendMessageResponse struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	Message    string   `json:"message,omitempty"`
	ProductIDs []string `json:"product_ids"`
}

type HeartbeatSubscriptionData struct {
	Feed string `json:"feed"`
	Time int    `json:"time"`
}

type CandlesTradeData struct {
	Feed      string `json:"feed"`
	Candle    Candle `json:"candle,omitempty"`
	ProductID string `json:"product_id"`
}

// -------------------------------------------------------------------------------------- //

type Candle struct {
	Time   int    `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume int    `json:"volume"`
}
