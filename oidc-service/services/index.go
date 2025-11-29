package services

import (
	"local/bomboclat-oidc-service/services/sessions"
	"local/bomboclat-oidc-service/services/users"
	custom_types "local/bomboclat-oidc-service/types"

	"github.com/redis/go-redis/v9"
)

var (
	UserService    users.UserService
	SessionService sessions.SessionService
)

func InjectDBToServices(db *custom_types.Postgres) {
	UserService.DBConn = db
	SessionService.DBConn = db
}

func InjectRedisClientToServices(c *redis.Client) {
	UserService.RedisClient = c
	SessionService.RedisClient = c
}
