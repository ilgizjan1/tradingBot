package postgresRepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"trade-bot/internal/pkg/models"
)

var (
	ErrCouldNotRollbackTransaction = errors.New("could not rollback transaction")
	ErrGetUsersOrder               = errors.New("get user orders")
)

type KrakenOrdersManagerPostgres struct {
	db *sqlx.DB
}

func NewKrakenOrdersManagerPostgres(db *sqlx.DB) *KrakenOrdersManagerPostgres {
	return &KrakenOrdersManagerPostgres{db: db}
}

const createOrderQuery = `
	INSERT INTO orders(order_id, user_id, cli_order_id, type, symbol, quantity, side, filled,
	                  timestamp, last_update_timestamp, price)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8,
	                  $9, $10, $11)`

const createUsersOrdersQuery = `
	INSERT INTO users_orders(user_id, order_id) VALUES ($1, $2)
`

func (k *KrakenOrdersManagerPostgres) CreateOrder(userID int, order models.Order) error {
	tx, err := k.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createOrderQuery, order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
		order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return ErrCouldNotRollbackTransaction
		}
		return err
	}

	if _, err = tx.Exec(createUsersOrdersQuery, userID, order.ID); err != nil {
		if err := tx.Rollback(); err != nil {
			return ErrCouldNotRollbackTransaction
		}
		return err
	}

	return tx.Commit()
}

const getOrderByIDQuery = `
SELECT * FROM orders WHERE order_id like $1
`

func (k *KrakenOrdersManagerPostgres) GetOrder(orderID string) (models.Order, error) {
	var order models.Order
	err := k.db.Get(&order, getOrderByIDQuery, orderID)
	return order, err
}

const getUserOrdersQuery = `SELECT * FROM orders WHERE user_id=$1`

func (k *KrakenOrdersManagerPostgres) GetUserOrders(userID int) ([]models.Order, error) {
	rows, err := k.db.Query(getUserOrdersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrGetUsersOrder, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order

		if err := rows.Scan(&order.ID, &order.UserID, &order.ClientOrderID, &order.Type, &order.Symbol, &order.Quantity,
			&order.Side, &order.Filled, &order.Timestamp, &order.LastUpdateTimestamp, &order.Price); err != nil {
			return nil, fmt.Errorf("%s: %w", ErrGetUsersOrder, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
