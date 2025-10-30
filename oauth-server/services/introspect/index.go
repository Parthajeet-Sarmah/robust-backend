package introspect

import (
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/redis/go-redis/v9"
)

type IntrospectService struct {
	RedisClient *redis.Client
	DBConn      *custom_types.Postgres
}
