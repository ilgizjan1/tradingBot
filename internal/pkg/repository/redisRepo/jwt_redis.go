package redisRepo

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"trade-bot/pkg/utils"
)

type JWTRedis struct {
	client *redis.Client
}

func NewJWTRedis(client *redis.Client) *JWTRedis {
	return &JWTRedis{client: client}
}

func (r *JWTRedis) CreateJWT(userID int, td utils.TokenDetails) (string, error) {
	at := time.Unix(td.AtExpires, 0)
	now := time.Now()

	errAccess := r.client.Set(context.Background(), td.AccessUUID, strconv.Itoa(userID), at.Sub(now)).Err()
	if errAccess != nil {
		return "", errAccess
	}

	return td.AccessToken, nil
}

func (r *JWTRedis) GetJWTUserID(ad utils.AccessDetails) (int, error) {
	strUserID, err := r.client.Get(context.Background(), ad.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *JWTRedis) DeleteJWT(ad utils.AccessDetails) error {
	_, err := r.client.Del(context.Background(), ad.AccessUUID).Result()
	return err
}
