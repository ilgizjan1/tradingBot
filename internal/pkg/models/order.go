package models

type Order struct {
	ID                  string  `json:"id" db:"order_id"`
	UserID              int     `json:"user_id" db:"user_id"`
	ClientOrderID       string  `json:"client_order_id" db:"cli_order_id"`
	Type                string  `json:"type" db:"type"`
	Symbol              string  `json:"symbol" db:"symbol"`
	Quantity            float64 `json:"quantity" db:"quantity"`
	Side                string  `json:"side" db:"side"`
	Filled              float64 `json:"filled" db:"filled"`
	Timestamp           string  `json:"timestamp" db:"timestamp"`
	LastUpdateTimestamp string  `json:"last_update_timestamp" db:"last_update_timestamp"`
	Price               float64 `json:"price" db:"price"`
}
