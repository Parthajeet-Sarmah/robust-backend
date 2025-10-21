package services

import (
	"local/bomboclat-oauth-server/services/authorization"
	"local/bomboclat-oauth-server/services/clients"
	"local/bomboclat-oauth-server/services/users"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/redis/go-redis/v9"
)

var (
	AuthorizationService authorization.AuthorizationService
	UserService          users.UserService
	ClientService        clients.ClientService
)

func InjectDBToServices(db *custom_types.Postgres) {
	UserService.DBConn = db
	AuthorizationService.DBConn = db
	ClientService.DBConn = db
}

func InjectRedisClientToServices(c *redis.Client) {
	AuthorizationService.RedisClient = c
	UserService.RedisClient = c
	ClientService.RedisClient = c
}
