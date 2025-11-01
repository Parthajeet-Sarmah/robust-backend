package users

import (
	custom_types "local/bomboclat-oidc-service/types"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	RedisClient *redis.Client
	DBConn      *custom_types.Postgres
}
