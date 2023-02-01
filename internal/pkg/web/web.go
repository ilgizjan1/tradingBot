package web

import (
	"context"

	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/web/webKraken"
	"trade-bot/pkg/krakenFuturesSDK"
	"trade-bot/pkg/krakenFuturesWSSDK"
)

type KrakenOrdersManager interface {
	SendOrder(args krakenFuturesSDK.SendOrderArguments) (krakenFuturesSDK.SendStatus, error)
	EditOrder(args krakenFuturesSDK.EditOrderArguments) (krakenFuturesSDK.EditStatus, error)
	CancelOrder(args krakenFuturesSDK.CancelOrderArguments) (krakenFuturesSDK.CancelStatus, error)
	CancelAllOrders(symbol string) (krakenFuturesSDK.CancelAllStatus, error)
	ParseSendStatusToExecutedOrder(userID int, sendStatus krakenFuturesSDK.SendStatus) (models.Order, error)
}

type KrakenAnalyzer interface {
	LookForCandles(ctx context.Context, feed string, productsIDs []string) (<-chan krakenFuturesWSSDK.Candle, error)
}

type Web struct {
	KrakenOrdersManager
	KrakenAnalyzer
}

func NewWeb(krakenAPISDK *krakenFuturesSDK.API, krakenWebsocketSDK *krakenFuturesWSSDK.WSAPI) *Web {
	return &Web{
		KrakenOrdersManager: webKraken.NewKrakenOrdersManagerWebSDK(krakenAPISDK),
		KrakenAnalyzer:      webKraken.NewKrakenAnalyzerWebSDK(krakenWebsocketSDK),
	}
}
