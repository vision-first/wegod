package facades

import (
	"sync"
	"wegod/internal/pkg/config"
	"github.com/995933447/log-go"
	"wegod/internal/pkg/redis"
)

var (
	newRedisGroupMu sync.Mutex
	redisGroup *redis.RedisGroup
)

func RedisGroup(logger *log.Logger) *redis.RedisGroup {
	if redisGroup == nil {
		newRedisGroupMu.Lock()
		defer newRedisGroupMu.Unlock()
		var redisNodes []*redis.RedisNode
		for _, nodeConfig := range config.Conf.Redis.Nodes {
			redisNodes = append(redisNodes, redis.NewRedisNode(nodeConfig.Host, nodeConfig.Port, nodeConfig.Password))
		}
		redisGroup = redis.NewRedisGroup(redisNodes, logger)
	}

	return redisGroup
}
