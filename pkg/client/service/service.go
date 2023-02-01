package service

import (
	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"
)

type Authorization interface {
	SignUp(input models.SignUpInput) (models.SignUpResponse, error)
	SignIn(input models.SignInInput) (models.SignInResponse, error)
	Logout(input models.LogoutInput) (models.LogoutResponse, error)
}

type OrdersManager interface {
	SendOrder(input models.SendOrderInput) (models.SendOrderResponse, error)
	StartTrading(input models.StartTradingInput) (<-chan *models.StartTradingResponse, <-chan error, error)
	GetUserOrders(input models.GetUserOrdersInput) (models.GetUserOrdersResponse, error)
}

type Service struct {
	Authorization
	OrdersManager
}

func NewService(client app.ClientActions) *Service {
	return &Service{
		Authorization: NewAuthService(client),
		OrdersManager: NewOrdersManagerService(client),
	}
}
