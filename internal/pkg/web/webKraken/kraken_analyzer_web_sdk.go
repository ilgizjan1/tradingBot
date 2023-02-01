package webKraken

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"trade-bot/pkg/krakenFuturesWSSDK"
)

var (
	ErrConvertTradeDataToCandle = errors.New("convert trade data to candle")
	ErrLookForCandles           = errors.New("look for candles")
)

const unixTimeLen = 10

type KrakenAnalyzerWebSDK struct {
	krakenWebsocketAPI *krakenFuturesWSSDK.WSAPI
}

func NewKrakenAnalyzerWebSDK(krakenWebsocketAPI *krakenFuturesWSSDK.WSAPI) *KrakenAnalyzerWebSDK {
	return &KrakenAnalyzerWebSDK{krakenWebsocketAPI: krakenWebsocketAPI}
}

func (k *KrakenAnalyzerWebSDK) LookForCandles(ctx context.Context, feed string, productsIDs []string) (<-chan krakenFuturesWSSDK.Candle, error) {
	tradeDataCh, err := k.krakenWebsocketAPI.CandlesTrade(ctx, feed, productsIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrLookForCandles, err)
	}

	candleCh, errCh := convertTradeDataToCandle(tradeDataCh)
	go logErrors(errCh)

	filteredCandles := filterCandles(candleCh)

	filteredUnixTimeCandles, errCh := filterCandlesUnixTime(filteredCandles)
	go logErrors(errCh)

	return filteredUnixTimeCandles, nil
}

func logErrors(errs <-chan error) {
	for err := range errs {
		log.Warn(err)
	}
}

func convertTradeDataToCandle(tradeData <-chan *krakenFuturesWSSDK.CandlesTradeData) (<-chan krakenFuturesWSSDK.Candle, <-chan error) {
	errCh := make(chan error, 1)
	candlesChan := make(chan krakenFuturesWSSDK.Candle)

	go func() {
		defer close(candlesChan)
		defer close(errCh)

		for data := range tradeData {
			if data.Feed == "error" {
				errCh <- fmt.Errorf("%s: error feed sended", ErrConvertTradeDataToCandle)
				continue
			}

			candlesChan <- data.Candle
		}
	}()

	return candlesChan, errCh
}

func filterCandles(candles <-chan krakenFuturesWSSDK.Candle) <-chan krakenFuturesWSSDK.Candle {
	candlesChan := make(chan krakenFuturesWSSDK.Candle)

	go func() {
		defer close(candlesChan)

		var lastUpdateTime *int

		for candle := range candles {
			if lastUpdateTime == nil {
				t := candle.Time
				lastUpdateTime = &t
				candlesChan <- candle
				continue
			}
			if candle.Time <= *lastUpdateTime {
				continue
			}

			*lastUpdateTime = candle.Time
			candlesChan <- candle
		}
	}()

	return candlesChan
}

func filterCandlesUnixTime(candles <-chan krakenFuturesWSSDK.Candle) (<-chan krakenFuturesWSSDK.Candle, <-chan error) {
	errCh := make(chan error, 1)
	candlesChan := make(chan krakenFuturesWSSDK.Candle)

	go func() {
		defer close(errCh)
		defer close(candlesChan)

		for candle := range candles {
			strCandleTime := strconv.Itoa(candle.Time)
			if len(strCandleTime) > unixTimeLen {
				lenDiff := len(strCandleTime) - unixTimeLen
				correctedTime := strCandleTime[:len(strCandleTime)-lenDiff]

				newTime, err := strconv.ParseInt(correctedTime, 10, 64)
				if err != nil {
					errCh <- err
					continue
				}

				candle.Time = int(newTime)
				candlesChan <- candle
			}
		}
	}()

	return candlesChan, errCh
}
