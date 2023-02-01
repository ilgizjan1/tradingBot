package webKraken

import (
	"fmt"

	"github.com/pkg/errors"

	"trade-bot/internal/pkg/models"
	"trade-bot/pkg/krakenFuturesSDK"
)

var (
	ErrSendOrder             = errors.New("web sdk: send order")
	ErrEditOrder             = errors.New("web sdk: edit order")
	ErrCancelOrder           = errors.New("web sdk: cancel order")
	ErrCancelAllOrders       = errors.New("web sdk: cancel all orders")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrUnknownSendStatusType = errors.New("unknown send status type")
)

type KrakenOrdersManagerWebSDK struct {
	api *krakenFuturesSDK.API
}

func NewKrakenOrdersManagerWebSDK(api *krakenFuturesSDK.API) *KrakenOrdersManagerWebSDK {
	return &KrakenOrdersManagerWebSDK{api: api}
}

func (k *KrakenOrdersManagerWebSDK) SendOrder(args krakenFuturesSDK.SendOrderArguments) (krakenFuturesSDK.SendStatus, error) {
	response, err := k.api.SendOrder(args)
	if err != nil {
		return krakenFuturesSDK.SendStatus{}, fmt.Errorf("%s: %w", ErrSendOrder, err)
	}

	if response.Error != "" {
		err := fmt.Errorf("err: %s, server time: %s, result: %s", response.Error, response.ServerTime, response.Result)
		return krakenFuturesSDK.SendStatus{}, fmt.Errorf("%s: %w", ErrSendOrder, err)
	}

	if !response.SendStatus.Status.IsSuccessStatus() {
		err := fmt.Errorf("%s: status: %s", ErrInvalidStatus, response.SendStatus.Status)
		return krakenFuturesSDK.SendStatus{}, fmt.Errorf("%s: %w", ErrSendOrder, err)
	}

	return response.SendStatus, nil
}

func (k *KrakenOrdersManagerWebSDK) EditOrder(args krakenFuturesSDK.EditOrderArguments) (krakenFuturesSDK.EditStatus, error) {
	response, err := k.api.EditOrder(args)
	if err != nil {
		return krakenFuturesSDK.EditStatus{}, fmt.Errorf("%s: %w", ErrEditOrder, err)
	}

	if response.Error != "" {
		err := fmt.Errorf("err: %s, server time: %s, result: %s", response.Error, response.ServerTime, response.Result)
		return krakenFuturesSDK.EditStatus{}, fmt.Errorf("%s: %w", ErrEditOrder, err)
	}

	if !response.EditStatus.Status.IsSuccessStatus() {
		err := fmt.Errorf("%s: status: %s", ErrInvalidStatus, response.EditStatus.Status)
		return krakenFuturesSDK.EditStatus{}, fmt.Errorf("%s: %w", ErrEditOrder, err)
	}

	return response.EditStatus, nil
}

func (k *KrakenOrdersManagerWebSDK) CancelOrder(args krakenFuturesSDK.CancelOrderArguments) (krakenFuturesSDK.CancelStatus, error) {
	response, err := k.api.CancelOrder(args)
	if err != nil {
		return krakenFuturesSDK.CancelStatus{}, fmt.Errorf("%s: %w", ErrCancelOrder, err)
	}

	if response.Error != "" {
		err := fmt.Errorf("err: %s, server time: %s, result: %s", response.Error, response.ServerTime, response.Result)
		return krakenFuturesSDK.CancelStatus{}, fmt.Errorf("%s: %w", ErrCancelOrder, err)
	}

	if !response.CancelStatus.Status.IsSuccessStatus() {
		err := fmt.Errorf("%s: status: %s", ErrInvalidStatus, response.CancelStatus.Status)
		return krakenFuturesSDK.CancelStatus{}, fmt.Errorf("%s: %w", ErrCancelOrder, err)
	}

	return response.CancelStatus, nil
}

func (k *KrakenOrdersManagerWebSDK) CancelAllOrders(symbol string) (krakenFuturesSDK.CancelAllStatus, error) {
	response, err := k.api.CancelAllOrders(symbol)
	if err != nil {
		return krakenFuturesSDK.CancelAllStatus{}, fmt.Errorf("%s: %w", ErrCancelAllOrders, err)
	}

	if response.Error != "" {
		err := fmt.Errorf("err: %s, server time: %s, result: %s", response.Error, response.ServerTime, response.Result)
		return krakenFuturesSDK.CancelAllStatus{}, fmt.Errorf("%s: %w", ErrCancelAllOrders, err)
	}

	if !response.CancelStatus.Status.IsSuccessStatus() {
		err := fmt.Errorf("%s: status: %s", ErrInvalidStatus, response.CancelStatus.Status)
		return krakenFuturesSDK.CancelAllStatus{}, fmt.Errorf("%s: %w", ErrCancelAllOrders, err)
	}

	return response.CancelStatus, nil
}

func (k *KrakenOrdersManagerWebSDK) ParseSendStatusToExecutedOrder(userID int, sendStatus krakenFuturesSDK.SendStatus) (models.Order, error) {
	orderEvent := sendStatus.OrderEvents[0]

	if orderEvent.Type == "EXECUTION" {
		return models.Order{
			ID:                  orderEvent.OrderPriorExecution.OrderID,
			UserID:              userID,
			ClientOrderID:       orderEvent.OrderPriorExecution.CliOrderID,
			Type:                orderEvent.Type,
			Symbol:              orderEvent.OrderPriorExecution.Symbol,
			Quantity:            orderEvent.OrderPriorExecution.Quantity,
			Side:                orderEvent.OrderPriorExecution.Side,
			Price:               orderEvent.Price,
			Filled:              orderEvent.OrderPriorExecution.Filled,
			Timestamp:           orderEvent.OrderPriorExecution.Timestamp,
			LastUpdateTimestamp: orderEvent.OrderPriorExecution.LastUpdateTimestamp,
		}, nil
	}

	return models.Order{}, ErrUnknownSendStatusType
}
