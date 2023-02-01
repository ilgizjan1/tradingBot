package handler

import (
	"fmt"
	"net/http"
	"trade-bot/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrUserIdentity   = errors.New("user identity")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrUserNotFound   = errors.New("user not found")
	ErrAPIKeyNotFound = errors.New("api key not found")
	ErrInvalidAPIKey  = errors.New("invalid api key")
)

const (
	userIDCtx            = "userID"
	userPublicAPIKeyCtx  = "publicAPIKey"
	userPrivateAPIKeyCtx = "privateAPIKey"
)

func (h *Handler) userIdentity(c *gin.Context) {
	bearerToken, err := utils.GetBearerToken(c.Request)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized,
			fmt.Sprintf("%s: %s", ErrUserIdentity.Error(), err.Error()))
		return
	}

	userID, err := h.services.Authorization.GetUserIDByJWT(bearerToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized,
			fmt.Sprintf("%s: %s", ErrUserIdentity.Error(), err.Error()))
		return
	}

	publicKey, privateKey, err := h.services.Authorization.GetUserAPIKeys(userID)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized,
			fmt.Sprintf("%s: %s", ErrUserIdentity.Error(), err.Error()))
		return
	}

	c.Set(userIDCtx, userID)
	c.Set(userPublicAPIKeyCtx, publicKey)
	c.Set(userPrivateAPIKeyCtx, privateKey)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userIDCtx)
	if !ok {
		return 0, ErrUserNotFound
	}

	intID, ok := id.(int)
	if !ok {
		return 0, ErrInvalidUserID
	}
	return intID, nil
}

func getUserAPIKeys(c *gin.Context) (string, string, error) {
	public, err := getUserAPIKey(c, userPublicAPIKeyCtx)
	if err != nil {
		return "", "", err
	}
	private, err := getUserAPIKey(c, userPrivateAPIKeyCtx)
	if err != nil {
		return "", "", err
	}
	return public, private, nil
}

func getUserAPIKey(c *gin.Context, keyType string) (string, error) {
	value, ok := c.Get(keyType)
	if !ok {
		return "", fmt.Errorf("%s: %s", ErrAPIKeyNotFound, keyType)
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s: %s", ErrInvalidAPIKey, value)
	}
	return stringValue, nil
}
