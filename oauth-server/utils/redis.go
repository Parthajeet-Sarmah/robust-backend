package utils

import (
	"context"
	"fmt"
	custom_errors "local/bomboclat-oauth-server/errors"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	if client == nil {
		return nil, custom_errors.RedisCouldNotCreateClientError(nil)
	}

	return client, nil
}

func GetValueFromHash(c *redis.Client, hash string) (map[string]string, error) {
	ctx := context.Background()

	res, err := c.HGetAll(ctx, hash).Result()

	if err != nil {
		return nil, custom_errors.RedisGetHashError(err)
	}

	if res == nil {
		return nil, custom_errors.RedisGetHashNoResourceFoundError(fmt.Errorf("No resource found for the given hash: %s", hash))
	}

	return res, nil
}

func SetValueToHash(c *redis.Client, hash string, resource map[string]string) error {
	ctx := context.Background()

	_, err := c.HSet(ctx, hash, resource).Result()

	if err != nil {
		return custom_errors.RedisSetHasError(fmt.Errorf("Error while setting for hash %s : %s", hash, err.Error()))
	}

	return nil

}
