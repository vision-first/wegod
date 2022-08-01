package facades

import (
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/redis"
)

func RedisGroup() *redis.RedisGroup {
	return facades.RedisGroup(MustLogger())
}
