package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type websocketErrResponse struct {
	Message string `json:"message"`
}

func newWebsocketErrResponse(c *gin.Context, code int, ws *websocket.Conn, message string) {
	if err := ws.WriteJSON(websocketErrResponse{Message: message}); err != nil {
		log.Error(err.Error())
	}
	c.AbortWithStatus(code)
	log.Error(message)
}
