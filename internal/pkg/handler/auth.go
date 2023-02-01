package handler

import (
	"net/http"
	"trade-bot/internal/pkg/models"
	"trade-bot/pkg/utils"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInputBody = "invalid input body"
)

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "access_token"
// @Failure 400,404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, ErrInvalidInputBody)
		return
	}

	accessToken, err := h.services.Authorization.GenerateJWT(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, ErrInvalidInputBody)
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Logout
// @Security ApiKeyAuth
// @Tags auth
// @Description logout account
// @ID logout-account
// @Produce  json
// @Success 200 {string} string "message"
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure default {object} errResponse
// @Router /auth/logout [delete]
func (h *Handler) logout(c *gin.Context) {
	token, err := utils.GetBearerToken(c.Request)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.services.Authorization.LogoutUser(token); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully logged out",
	})
}
