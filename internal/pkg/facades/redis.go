package facades

import (
	"github.com/995933447/log-go"
	"github.com/995933447/redisgroup"
	"github.com/vision-first/wegod/internal/pkg/config"
	"sync"
)

var (
	newRedisGroupMu sync.Mutex
	redisGroup *redisgroup.Group
)

func RedisGroup(logger *log.Logger) *redisgroup.Group {
	if redisGroup != nil {
		return redisGroup
	}

	newRedisGroupMu.Lock()
	defer newRedisGroupMu.Unlock()

	if redisGroup != nil {
		return redisGroup
	}

	var redisNodes []*redisgroup.Node
	for _, nodeConfig := range config.Conf.Redis.Nodes {
		redisNodes = append(redisNodes, redisgroup.NewNode(nodeConfig.Host, nodeConfig.Port, nodeConfig.Password))
	}
	redisGroup = redisgroup.NewGroup(redisNodes, logger)

	return redisGroup
}
