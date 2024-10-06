package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "redisCache:6379",
		Password: password,
		DB:       0,
	})

}

func SetCode(code string, user_id uint) (string, error) {

	ctx := context.Background()

	userID := fmt.Sprintf("%d", user_id)

	err := RedisClient.Set(ctx, userID, code, 30*time.Minute).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetCode(user_id uint) (string, error) {
	ctx := context.Background()

	userIDStr := fmt.Sprintf("%d", user_id)

	cnt := RedisClient.Get(ctx, userIDStr).String()

	return cnt, nil
}
