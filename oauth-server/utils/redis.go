package utils

import (
	"context"
	"fmt"
	custom_errors "local/bomboclat-oauth-server/errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type HashOptions struct {
	Expiry time.Duration
}

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

	if len(res) == 0 {
		return nil, custom_errors.RedisGetHashNoResourceFoundError(fmt.Errorf("No resource found for the given hash: %s", hash))
	}

	return res, nil
}

func SetValueToHash(c *redis.Client, hash string, resource map[string]string, expiry ...time.Duration) error {
	ctx := context.Background()

	_, err := c.HSet(ctx, hash, resource).Result()

	//Check if any expiry is provided, if yes add it to redis
	if len(expiry) > 1 {
		return custom_errors.RedisTooManyExpiryTimesError(fmt.Errorf("Too many expiry times while setting hash: %s", hash))
	} else if len(expiry) == 1 {
		c.Expire(ctx, hash, expiry[0])
	}

	if err != nil {
		return custom_errors.RedisSetHasError(fmt.Errorf("Error while setting for hash %s : %s", hash, err.Error()))
	}

	return nil

}
