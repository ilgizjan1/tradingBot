package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/pkg/krakenFuturesSDK"
)

// @Summary SendOrder
// @Security ApiKeyAuth
// @Tags orderManager
// @Description sendOrder to kraken futures API
// @ID sendOrder
// @Accept  json
// @Produce  json
// @Param input body krakenFuturesSDK.SendOrderArguments true "send order info"
// @Success 200 {string} string "order_id"
// @Failure 400,401,404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /orderManager/send-order [post]
func (h *Handler) sendOrder(c *gin.Context) {
	var input krakenFuturesSDK.SendOrderArguments

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	order, err := h.services.KrakenOrdersManager.SendOrder(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

type tradingDetails struct {
	Event          string               `json:"event"`
	TradingDetails types.TradingDetails `json:"trading_details,omitempty"`
}

const cancelEvent = "cancel_trading"
const startTrading = "start_trading"

func (h *Handler) startTrade(c *gin.Context) {
	conn, err := h.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	defer conn.Close()

	userID, err := getUserID(c)
	if err != nil {
		newWebsocketErrResponse(c, http.StatusUnauthorized, conn, err.Error())
		return
	}

	var input tradingDetails
	if err := conn.ReadJSON(&input); err != nil {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}
	if err := h.validate.Struct(input); err != nil {
		newWebsocketErrResponse(c, http.StatusBadRequest, conn, err.Error())
		return
	}
	if input.Event != startTrading {
		newWebsocketErrResponse(c, http.StatusBadRequest, conn, err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	var isCancelled bool
	defer cancel()

	go func() {
		defer cancel()

		var tradingDetails tradingDetails
		for {
			if err := conn.ReadJSON(&tradingDetails); err != nil {
				return
			}
			if tradingDetails.Event == cancelEvent {
				isCancelled = true
				return
			}
		}
	}()

	order, err := h.services.KrakenOrdersManager.StartTrading(ctx, userID, input.TradingDetails)
	if err != nil && !isCancelled {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}

	if isCancelled {
		err := conn.WriteJSON(struct {
			Message string `json:"message"`
		}{Message: "trading have been canceled"})

		if err != nil {
			newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
			return
		}
		return
	}

	if err := conn.WriteJSON(order); err != nil {
		newWebsocketErrResponse(c, http.StatusInternalServerError, conn, err.Error())
		return
	}
}

// @Summary MyOrders
// @Security ApiKeyAuth
// @Tags orderManager
// @Description get all orders of user
// @ID myOrders
// @Produce  json
// @Success 200 {object} []models.Order
// @Failure 401,404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /orderManager/my-orders [get]
func (h *Handler) myOrders(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	orders, err := h.services.KrakenOrdersManager.GetUserOrders(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}
