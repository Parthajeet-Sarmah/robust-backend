package utils

import (
	"context"
	"log"

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
		return nil, &RedisCouldNotCreateClient{}
	}

	return client, nil
}

func GetValueFromHash(c *redis.Client, hash string) (map[string]string, error) {
	ctx := context.Background()

	res, err := c.HGetAll(ctx, hash).Result()

	if err != nil {
		log.Print("Error while getting hash resource from Redis!")
		return nil, &RedisGetHashError{}
	}

	if res == nil {
		log.Print("No resource found for the given hash!")
		return nil, &RedisGetHashNoResourceFoundError{}
	}

	return res, nil
}

func SetValueToHash(c *redis.Client, hash string, resource map[string]string) error {
	ctx := context.Background()

	_, err := c.HSet(ctx, hash, resource).Result()

	if err != nil {
		log.Print("Error while setting hash resource to Redis!")
		return &RedisSetHasError{}
	}

	return nil

}
