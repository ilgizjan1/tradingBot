package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type errResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error(message)
	c.AbortWithStatusJSON(statusCode, errResponse{Message: message})
}
