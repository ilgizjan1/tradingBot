package redisRepo

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"trade-bot/configs"
)

var (
	ErrNewRedisDB              = errors.New("new redis db")
	ErrUnableToLocateRedisPort = errors.New("unable to locate redis port")
)
var redisContext = context.Background()

func NewRedisClient(cfg configs.RedisDatabaseConfiguration) (*redis.Client, error) {
	addr := cfg.Host + ":" + cfg.Port
	if len(addr) == 0 {
		return nil, fmt.Errorf("%s: %w", ErrNewRedisDB, ErrUnableToLocateRedisPort)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(redisContext).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrNewRedisDB, err)
	}
	return client, nil
}
