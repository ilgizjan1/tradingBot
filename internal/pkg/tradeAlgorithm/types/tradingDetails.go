package types

type TradingDetails struct {
	OrderType        string  `json:"order_type" validate:"required"`
	Symbol           string  `json:"symbol" validate:"required"`
	Side             string  `json:"side" validate:"required"`
	Size             uint    `json:"size" validate:"required,gte=0"`
	StopLossBorder   float64 `json:"stop_loss_border" validate:"required,gte=0"`
	TakeProfitBorder float64 `json:"take_profit_border" validate:"required,gte=0"`
	BuyPrice         float64
}
