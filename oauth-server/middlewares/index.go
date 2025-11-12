package middlewares

import (
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/redis/go-redis/v9"
)

type MiddlewareServiceContainer struct {
	DBConn      *custom_types.Postgres
	RedisClient *redis.Client
}

var (
	MiddlewareService MiddlewareServiceContainer
)

func InjectDBToServices(db *custom_types.Postgres) {
	MiddlewareService.DBConn = db
}

func InjectRedisClientToServices(c *redis.Client) {
	MiddlewareService.RedisClient = c
}
