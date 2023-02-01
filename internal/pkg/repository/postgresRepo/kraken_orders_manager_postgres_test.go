package postgresRepo

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"trade-bot/internal/pkg/models"
)

func TestKrakenOrdersManagerPostgres_CreateOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewKrakenOrdersManagerPostgres(sqlxDB)

	type args struct {
		order  models.Order
		userID int
	}
	type mockBehaviour func(userID int, order models.Order)

	tests := []struct {
		name    string
		input   args
		mock    mockBehaviour
		wantErr bool
	}{
		{
			name: "OK",
			input: args{
				order: models.Order{
					ID:                  "1",
					UserID:              1,
					ClientOrderID:       "1",
					Type:                "type",
					Symbol:              "symbol",
					Quantity:            10,
					Side:                "buy",
					Filled:              2,
					Timestamp:           "timestamp",
					LastUpdateTimestamp: "timestamp",
					Price:               10,
				},
				userID: 1,
			},
			mock: func(userID int, order models.Order) {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO orders").
					WithArgs(order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
						order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO users_orders").WithArgs(userID, order.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Empty fields in 1 insert",
			input: args{
				order: models.Order{
					ID:     "",
					UserID: 1,
				},
				userID: 1,
			},
			mock: func(userID int, order models.Order) {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO orders").
					WithArgs(order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
						order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty fields in 2 insert",
			input: args{
				order: models.Order{
					ID:                  "1",
					UserID:              1,
					ClientOrderID:       "1",
					Type:                "type",
					Symbol:              "symbol",
					Quantity:            10,
					Side:                "buy",
					Filled:              2,
					Timestamp:           "timestamp",
					LastUpdateTimestamp: "timestamp",
					Price:               10,
				},
			},
			mock: func(userID int, order models.Order) {
				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO orders").
					WithArgs(order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
						order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO users_orders").WithArgs(userID, order.ID).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input.userID, test.input.order)

			err := r.CreateOrder(test.input.userID, test.input.order)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestKrakenOrdersManagerPostgres_GetOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewKrakenOrdersManagerPostgres(sqlxDB)

	type args struct {
		inputOrderID string
		order        models.Order
	}
	type mockBehaviour func(orderID string, order models.Order)

	tests := []struct {
		name    string
		input   args
		want    models.Order
		mock    mockBehaviour
		wantErr bool
	}{
		{
			name: "OK",
			input: args{
				inputOrderID: "1",
				order: models.Order{
					ID: "1",
				},
			},
			want: models.Order{
				ID: "1",
			},
			mock: func(orderID string, order models.Order) {
				rows := sqlmock.NewRows([]string{"order_id", "user_id", "cli_order_id", "type", "symbol", "quantity",
					"side", "filled", "timestamp", "last_update_timestamp", "price"}).
					AddRow(order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
						order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price)
				mock.ExpectQuery("SELECT (.+) FROM orders").
					WithArgs(orderID).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Not found",
			input: args{
				inputOrderID: "1",
				order: models.Order{
					ID: "1",
				},
			},
			want: models.Order{
				ID: "1",
			},
			mock: func(orderID string, order models.Order) {
				rows := sqlmock.NewRows([]string{"order_id", "user_id", "cli_order_id", "type", "symbol", "quantity",
					"side", "filled", "timestamp", "last_update_timestamp", "price"})
				mock.ExpectQuery("SELECT (.+) FROM orders").
					WithArgs(orderID).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input.inputOrderID, test.input.order)

			gorOrder, err := r.GetOrder(test.input.inputOrderID)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, gorOrder)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestKrakenOrdersManagerPostgres_GetUserOrders(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewKrakenOrdersManagerPostgres(sqlxDB)

	tests := []struct {
		name    string
		userID  int
		order   models.Order
		mock    func(userID int, order models.Order)
		want    []models.Order
		wantErr bool
	}{
		{
			name:   "OK",
			userID: 1,
			order: models.Order{
				ID:                  "1",
				UserID:              1,
				ClientOrderID:       "1",
				Type:                "type",
				Symbol:              "symbol",
				Quantity:            10,
				Side:                "buy",
				Filled:              10,
				Timestamp:           "time",
				LastUpdateTimestamp: "time",
				Price:               100,
			},
			mock: func(userID int, order models.Order) {
				rows := sqlmock.NewRows([]string{"order_id", "user_id", "cli_order_id", "type", "symbol", "quantity",
					"side", "filled", "timestamp", "last_update_timestamp", "price"}).
					AddRow(order.ID, order.UserID, order.ClientOrderID, order.Type, order.Symbol, order.Quantity,
						order.Side, order.Filled, order.Timestamp, order.LastUpdateTimestamp, order.Price)
				mock.ExpectQuery("SELECT (.+) FROM orders").
					WithArgs(userID).WillReturnRows(rows)
			},
			want: []models.Order{{
				ID:                  "1",
				UserID:              1,
				ClientOrderID:       "1",
				Type:                "type",
				Symbol:              "symbol",
				Quantity:            10,
				Side:                "buy",
				Filled:              10,
				Timestamp:           "time",
				LastUpdateTimestamp: "time",
				Price:               100,
			}},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.userID, test.order)

			gorOrders, err := r.GetUserOrders(test.userID)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, gorOrders)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
