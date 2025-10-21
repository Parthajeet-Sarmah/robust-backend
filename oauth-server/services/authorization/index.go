package authorization

import (
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/redis/go-redis/v9"
)

type AuthorizationService struct {
	RedisClient *redis.Client
	DBConn      *custom_types.Postgres
}
