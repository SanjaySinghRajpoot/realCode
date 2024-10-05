package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: password,
		DB:       0,
	})

}

func SetIPAddress(IPaddr string, Count int) (string, error) {

	ctx := context.Background()

	err := RedisClient.Set(ctx, IPaddr, Count, 30*time.Minute).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetIPAddress(IPaddr string) (int, error) {

	ctx := context.Background()

	cnt, err := RedisClient.Get(ctx, IPaddr).Int()

	if err != nil {
		return -1, err
	}

	return cnt, nil
}

func CheckIPAddressKey(IPAddr string) (bool, error) {
	// Create a context
	ctx := context.Background()

	// Check if the key exists in Redis
	exists, err := RedisClient.Exists(ctx, IPAddr).Result()

	if err != nil {
		// Handle error if any
		return false, err
	}

	// Return true if the key exists, false otherwise
	return exists == 1, nil
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
