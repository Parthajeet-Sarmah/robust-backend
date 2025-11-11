package services

import (
	"local/bomboclat-oauth-server/services/authorization"
	"local/bomboclat-oauth-server/services/clients"
	"local/bomboclat-oauth-server/services/introspect"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/redis/go-redis/v9"
)

var (
	AuthorizationService authorization.AuthorizationService
	ClientService        clients.ClientService
	IntrospectService    introspect.IntrospectService
)

func InjectDBToServices(db *custom_types.Postgres) {
	AuthorizationService.DBConn = db
	ClientService.DBConn = db
	IntrospectService.DBConn = db
}

func InjectRedisClientToServices(c *redis.Client) {
	AuthorizationService.RedisClient = c
	ClientService.RedisClient = c
	IntrospectService.RedisClient = c
}
