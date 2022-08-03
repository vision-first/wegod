package facades

import (
	"github.com/995933447/redisgroup"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

func RedisGroup() *redisgroup.Group {
	return facades.RedisGroup(MustLogger())
}
