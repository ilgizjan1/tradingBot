package algorithms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesWSSDK"
)

var (
	ErrStartAnalyzing     = errors.New("start analyzing")
	ErrUnableToGetCandles = errors.New("unable to get candles")
)

type StopLossTakeProfitAlgo struct {
	krakenWebsocketSDK web.KrakenAnalyzer
}

func NewStopLossTakeProfitAlgo(krakenAnalyzer web.KrakenAnalyzer) *StopLossTakeProfitAlgo {
	return &StopLossTakeProfitAlgo{
		krakenWebsocketSDK: krakenAnalyzer,
	}
}

func (a *StopLossTakeProfitAlgo) StartAnalyzing(ctx context.Context, buyTime time.Time, details types.TradingDetails) error {
	candles, err := a.krakenWebsocketSDK.LookForCandles(ctx, krakenFuturesWSSDK.OneMinuteCandlesFeed, []string{details.Symbol})
	if err != nil {
		return fmt.Errorf("%s: %w", ErrStartAnalyzing, err)
	}

	for candle := range candles {
		price, err := strconv.ParseFloat(candle.Close, 64)

		if time.Unix(int64(candle.Time), 0).Before(buyTime) {
			continue
		}

		if err != nil {
			return fmt.Errorf("%s: %w", ErrStartAnalyzing, err)
		}

		if price > details.BuyPrice+details.TakeProfitBorder {
			return nil
		}
		if price < details.BuyPrice-details.StopLossBorder {
			return nil
		}
	}

	return fmt.Errorf("%s: %s", ErrStartAnalyzing, ErrUnableToGetCandles)
}
