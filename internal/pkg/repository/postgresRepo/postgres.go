package postgresRepo

import (
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // driver for sqlx
	"github.com/jmoiron/sqlx"

	"trade-bot/configs"
)

var (
	ErrNewPostgresDB = errors.New("new postgres db")
	ErrPingDB        = errors.New("ping db")
)

func NewPostgresDB(cfg configs.PostgreDatabaseConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrNewPostgresDB, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrNewPostgresDB, ErrPingDB, err)
	}

	return db, nil
}
