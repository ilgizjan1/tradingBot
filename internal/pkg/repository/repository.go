package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository/postgresRepo"
	"trade-bot/internal/pkg/repository/redisRepo"
	"trade-bot/pkg/utils"
)

type Authorization interface {
	CreateUser(models.User) (int, error)
	GetUser(username string) (models.User, error)
	GetUserAPIKeys(userID int) (string, string, error)
}

type JWT interface {
	CreateJWT(userID int, td utils.TokenDetails) (string, error)
	GetJWTUserID(ad utils.AccessDetails) (int, error)
	DeleteJWT(ad utils.AccessDetails) error
}

type KrakenOrdersManager interface {
	CreateOrder(userID int, order models.Order) error
	GetUserOrders(userID int) ([]models.Order, error)
	GetOrder(orderID string) (models.Order, error)
}

type Repository struct {
	Authorization
	JWT
	KrakenOrdersManager
}

func NewRepository(db *sqlx.DB, jwtDB *redis.Client) *Repository {
	return &Repository{
		Authorization:       postgresRepo.NewAuthPostgres(db),
		JWT:                 redisRepo.NewJWTRedis(jwtDB),
		KrakenOrdersManager: postgresRepo.NewKrakenOrdersManagerPostgres(db),
	}
}
