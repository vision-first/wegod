package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/995933447/log-go"
	"github.com/go-redis/redis/v8"
)

type RedisNode struct {
	client *redis.Client
	host   string
	port   int
	hash   uint32
}

func hashKey(s string) uint32 {
	f := fnv.New32a()
	_, _ = f.Write([]byte(s))
	return f.Sum32()
}

func NewRedisNode(host string, port int, password string) *RedisNode {
	n := new(RedisNode)
	n.host = host
	n.port = port
	n.hash = hashKey(fmt.Sprintf("%s:%d", host, port))
	n.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
	})
	return n
}

type redisNodeSorter struct {
	nodes []*RedisNode
}

func (s *redisNodeSorter) Len() int {
	return len(s.nodes)
}

func (s *redisNodeSorter) Less(i, j int) bool {
	return s.nodes[i].hash < s.nodes[j].hash
}

func (s *redisNodeSorter) Swap(i, j int) {
	tmp := s.nodes[i]
	s.nodes[i] = s.nodes[j]
	s.nodes[j] = tmp
}

type RedisGroup struct {
	nodes []*RedisNode
	mu       sync.RWMutex
	logger *log.Logger
}
func NewRedisGroup(nodes []*RedisNode, logger *log.Logger) *RedisGroup {
	group := &RedisGroup{
		nodes: nodes,
		logger: logger,
	}

	group.sortNodes()

	return group
}

func (g *RedisGroup) sortNodes() {
	i := &redisNodeSorter{nodes: g.nodes}
	sort.Sort(i)
	g.nodes = i.nodes
}

func (g *RedisGroup) FindNodeForKey(key string) *RedisNode {
	hash := hashKey(key)
	g.mu.RLock()
	defer g.mu.RUnlock()
	if len(g.nodes) == 0 {
		return nil
	}
	i := 0
	for i+1 < len(g.nodes) {
		if hash >= g.nodes[i].hash && hash < g.nodes[i+1].hash {
			return g.nodes[i]
		}
		i++
	}
	return g.nodes[len(g.nodes)-1]
}

func (g *RedisGroup) GetNodes() []*RedisNode {
	return g.nodes
}

func (g *RedisGroup) GetClient(node *RedisNode) *redis.Client {
	return node.client
}

func (g *RedisGroup) Set(ctx context.Context, key string, val []byte, exp time.Duration) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Set: key %s exp %v", key, exp)
	err := node.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}

	return nil
}

func (g *RedisGroup) SetUint64(ctx context.Context, key string, val uint64, exp time.Duration) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Set: key %s exp %v", key, exp)
	err := node.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}

	return nil
}

func (g *RedisGroup) SetRange(ctx context.Context, key string, offset int64, value string) error {
	g.logger.Debugf(ctx, "redis: SetRange: key %s offset", key, offset)
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	err := node.client.SetRange(ctx, key, offset, value).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	return nil
}

func (g *RedisGroup) SetByJson(ctx context.Context, key string, j interface{}, exp time.Duration) error {
	val, err := json.Marshal(j)
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	// 空串这里先不考虑
	if len(val) == 0 {
		return errors.New("unsupported empty value")
	}
	return g.Set(ctx, key, val, exp)
}

func (g *RedisGroup) HLen(ctx context.Context, key string) (uint32, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}

	g.logger.Debugf(ctx, "redis: HLen: key %s", key)
	v := node.client.HLen(ctx, key)
	err := v.Err()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}

	return uint32(v.Val()), nil
}

func (g *RedisGroup) HSetByJson(ctx context.Context, key, subKey string, j interface{}, exp time.Duration) error {
	val, err := json.Marshal(j)
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx,  "redis: HSetJson: key %s subKey %+v", key, subKey)
	err = node.client.HSet(ctx, key, subKey, string(val)).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	if exp > 0 {
		err = node.client.Expire(ctx, key, exp).Err()
		if err != nil {
			g.logger.Errorf(ctx, "err:%v", err)
			return err
		}
	}
	return nil
}

func (g *RedisGroup) HSetNXByJson(ctx context.Context, key, subKey string, j interface{}, exp time.Duration, setSuccess *bool) error {
	val, err := json.Marshal(j)
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HSetJsonNX: key %s subKey %+v exp %v", key, subKey, exp)
	res := node.client.HSetNX(ctx, key, subKey, string(val))
	err = res.Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	if exp > 0 {
		err = node.client.Expire(ctx, key, exp).Err()
		if err != nil {
			g.logger.Errorf(ctx, "err:%v", err)
			return err
		}
	}
	if setSuccess != nil {
		*setSuccess = res.Val()
	}
	return nil
}

func (g *RedisGroup) Get(ctx context.Context, key string) ([]byte, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}

	g.logger.Debugf(ctx, "redis: Get: key %s", key)
	val, err := node.client.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			g.logger.Errorf(ctx, "err:%s", err)
		}
		return nil, err
	}
	return val, nil
}

func (g *RedisGroup) GetJson(ctx context.Context, key string, j interface{}) error {
	val, err := g.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	err = json.Unmarshal(val, j)
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	return nil
}

func (g *RedisGroup) GetUint64(ctx context.Context, key string) (uint64, error) {
	val, err := g.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return 0, redis.Nil
		}
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}

	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}

	return uint64(i), nil
}

func (g *RedisGroup) GetInt64(ctx context.Context, key string) (int64, error) {
	val, err := g.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return 0, redis.Nil
		}
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}

	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}

	return i, nil
}

func (g *RedisGroup) GetInt64Default(ctx context.Context, key string, def int64) (int64, error) {
	val, err := g.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return def, nil
		}
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}

	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}

	return i, nil
}

func (g *RedisGroup) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	g.logger.Debugf(ctx, "redis: HGetAll: key %s", key)
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	return node.client.HGetAll(ctx, key).Result()
}

func (g *RedisGroup) HKeys(ctx context.Context, key string) ([]string, error) {
	g.logger.Debugf(ctx, "redis: HKeys: key %s", key)
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	return node.client.HKeys(ctx, key).Result()
}

func (g *RedisGroup) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	g.logger.Debugf(ctx, "redis: HKeys: key %s", key)
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, 0, errors.New("not found available redis node")
	}
	return node.client.HScan(ctx, key, cursor, match, count).Result()
}

func (g *RedisGroup) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	g.logger.Debugf(ctx, "redis: ZKeys: key %s", key)
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, 0, errors.New("not found available redis node")
	}
	return node.client.ZScan(ctx, key, cursor, match, count).Result()
}

func (g *RedisGroup) HMGetJson(ctx context.Context, key, subKey string, j interface{}) error {
	g.logger.Debugf(ctx, "redis: HMGetJson: key %s subKey %+v", key, subKey)
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	values, err := node.client.HMGet(ctx, key, subKey).Result()
	if err != nil {
		g.logger.Errorf(ctx, "redis HMGet err:%v", err)
		return err
	}
	if len(values) == 1 {
		v := values[0]
		if v != nil {
			var buf []byte
			if p, ok := v.(string); ok {
				buf = []byte(p)
			} else if p, ok := v.([]byte); ok {
				buf = p
			}
			if buf != nil {
				if len(buf) > 0 {
					err = json.Unmarshal(buf, j)
					if err != nil {
						g.logger.Errorf(ctx, "err:%s", err)
						return err
					}
				}
				return nil
			}
		}
	}
	return redis.Nil
}

func (g *RedisGroup) HDel(ctx context.Context, key string, subKey ...string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HDel: key %s subKey %+v", key, subKey)
	delNum, err := node.client.HDel(ctx, key, subKey...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}
	return delNum, nil
}

func (g *RedisGroup) Del(ctx context.Context, key string) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Del: key %s", key)
	err := node.client.Del(ctx, key).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	return nil
}

func (g *RedisGroup) ZAdd(ctx context.Context, key string, values ...*redis.Z) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZAdd: key %s", key)
	err := node.client.ZAdd(ctx, key, values...).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	return nil
}

func (g *RedisGroup) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZCount: key %s min %s - max %s", key, min, max)
	v := node.client.ZCount(ctx, key, min, max)
	err := v.Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return v.Val(), nil
}

func (g *RedisGroup) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRangeByScore: key %s opt %v", key, opt)
	members, err := node.client.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return members, nil
}

func (g *RedisGroup) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRangeByScoreWithScores: key %s opt %v", key, opt)
	resultList, err := node.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return resultList, nil
}

func (g *RedisGroup) ZIncrBy(ctx context.Context, key string, increment float64, member string) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZIncrBy: key %s increment %v member %s", key, increment, member)
	_, err := node.client.ZIncrBy(ctx, key, increment, member).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	return nil
}

func (g *RedisGroup) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRange: key %s start %v stop %v", key, start, stop)
	members, err := node.client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return members, nil
}

func (g *RedisGroup) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRangeWithScores: key %s start %v stop %v", key, start, stop)
	resultList, err := node.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return resultList, nil
}

func (g *RedisGroup) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRem: key %s", key)
	delNum, err := node.client.ZRem(ctx, key, members...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return delNum, nil
}

func (g *RedisGroup) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZRemRangeByScore: key %s, min: %s, max: %s", key, min, max)
	delNum, err := node.client.ZRemRangeByScore(ctx, key, min, max).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return delNum, nil
}

func (g *RedisGroup) ZCard(ctx context.Context, key string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZCard: key %s", key)
	num, err := node.client.ZCard(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return num, nil
}

func (g *RedisGroup) SAdd(ctx context.Context, key string, values ...interface{}) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SAdd: key %s", key)
	err := node.client.SAdd(ctx, key, values...).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	return nil
}

func (g *RedisGroup) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SRem: key %s", key)
	delNum, err := node.client.SRem(ctx, key, members...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return delNum, nil
}

func (g *RedisGroup) SCard(ctx context.Context, key string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SCard: key %s", key)
	num, err := node.client.SCard(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return num, nil
}

func (g *RedisGroup) SIsMember(ctx context.Context, key string, members interface{}) (bool, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return false, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SIsMember: key %s", key)
	ok, err := node.client.SIsMember(ctx, key, members).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return false, err
	}
	return ok, nil
}

func (g *RedisGroup) SMembers(ctx context.Context, key string) ([]string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SMembers: key %s", key)
	members, err := node.client.SMembers(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return members, nil
}

func (g *RedisGroup) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HIncrBy: key %s field %s incr %d", key, field, incr)
	n, err := node.client.HIncrBy(ctx, key, field, incr).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}
	return n, nil
}

func (g *RedisGroup) IncrBy(ctx context.Context, key string, incr int64) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: IncrBy: key %s incr %d", key, incr)
	n, err := node.client.IncrBy(ctx, key, incr).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}
	return n, nil
}

func (g *RedisGroup) DecrBy(ctx context.Context, key string, decr int64) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: DecrBy: key %s decr %d", key, decr)
	n, err := node.client.DecrBy(ctx, key, decr).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return 0, err
	}
	return n, nil
}

func (g *RedisGroup) HSet(ctx context.Context, key, field string, val interface{}) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HSet: key %s field %s", key, field)
	err := node.client.HSet(ctx, key, field, val).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}

	return nil
}

func (g *RedisGroup) HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HMSet: key %s fields %v", key, fields)
	err := node.client.HMSet(ctx, key, fields).Err()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}

	return nil
}

func (g *RedisGroup) HGet(ctx context.Context, key, subKey string) (string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HGet: key %s subKey %+v", key, subKey)
	val, err := node.client.HGet(ctx, key, subKey).Result()
	if err != nil {
		if err != redis.Nil {
			g.logger.Errorf(ctx, "err:%v", err)
		}
		return "", err
	}
	return val, nil
}

func (g *RedisGroup) HGetJson(ctx context.Context, key, subKey string, j interface{}) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HGetJson: key %s subKey %+v", key, subKey)
	val, err := node.client.HGet(ctx, key, subKey).Result()
	if err != nil {
		if err != redis.Nil {
			g.logger.Errorf(ctx, "err:%v", err)
		}
		return err
	}
	err = json.Unmarshal([]byte(val), j)
	if err != nil {
		g.logger.Errorf(ctx, "err:%s", err)
		return err
	}
	return nil
}

func (g *RedisGroup) Expire(ctx context.Context, key string, expiration time.Duration) error {
	node := g.FindNodeForKey(key)
	if node == nil {
		return errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Expire: key %s exp %+v", key, expiration)
	_, err := node.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return err
	}
	return nil
}

func (g *RedisGroup) Exists(ctx context.Context, key string) (bool, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return false, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Exists: key %s", key)
	val, err := node.client.Exists(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return false, err
	}
	if val == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (g *RedisGroup) HExists(ctx context.Context, key, field string) (bool, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return false, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: HExists: key %s field %s", key, field)
	exists, err := node.client.HExists(ctx, key, field).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return false, err
	}
	return exists, nil
}

func (g *RedisGroup) ScriptRun(ctx context.Context, lua string, keys []string, args ...interface{}) (interface{}, error) {
	node := g.FindNodeForKey(keys[0])
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	script := redis.NewScript(lua)
	g.logger.Debugf(ctx, "lua run:%s", lua)
	result, err := script.Run(ctx, node.client, keys, args...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return result, nil
}

func (g *RedisGroup) EvalSha(ctx context.Context, luaSha1 string, keys []string, args ...interface{}) (interface{}, error) {
	node := g.FindNodeForKey(keys[0])
	if node == nil {
		return nil, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "lua eval:%s", luaSha1)
	result, err := node.client.EvalSha(ctx, luaSha1, keys, args...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return nil, err
	}
	return result, nil
}

func (g *RedisGroup) ScriptLoad(ctx context.Context, luaScript string) (string, error) {
	node := g.FindNodeForKey("")
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "lua load:%s", luaScript)
	luaSha1, err := node.client.ScriptLoad(ctx, luaScript).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return "", err
	}
	return luaSha1, nil
}

func (g *RedisGroup) Incr(ctx context.Context, key string) (int64, error) {
	node := g.FindNodeForKey("")
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Incr: key %s", key)
	val, err := node.client.Incr(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return val, nil
}

func (g *RedisGroup) Decr(ctx context.Context, key string) (int64, error) {
	node := g.FindNodeForKey("")
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Decr: key %s", key)
	val, err := node.client.Decr(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return val, nil
}

func (g *RedisGroup) ExpireAt(ctx context.Context, key string, expiredAt time.Time) (bool, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return false, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ExpireAt: key %s exp %v", key, expiredAt)
	ok, err := node.client.ExpireAt(ctx, key, expiredAt).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return false, err
	}
	return ok, nil
}

func (g *RedisGroup) LPop(ctx context.Context, key string) (string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LPop: key %s", key)
	val, err := node.client.LPop(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return "", err
	}
	return val, nil
}

func (g *RedisGroup) RPop(ctx context.Context, key string) (string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: RPop: key %s", key)
	val, err := node.client.RPop(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return "", err
	}
	return val, nil
}

func (g *RedisGroup) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LPush: key %s", key)
	count, err := node.client.LPush(ctx, key, values...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return count, nil
}

func (g *RedisGroup) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: RPush: key %s", key)
	count, err := node.client.RPush(ctx, key, values...).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return count, nil
}

func (g *RedisGroup) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return []string{}, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LRange: key %s start %d stop %d", key, start, stop)
	result, err := node.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return []string{}, err
	}
	return result, nil
}

func (g *RedisGroup) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LTrim: key %s start %d stop %d", key, start, stop)
	result, err := node.client.LTrim(ctx, key, start, stop).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return "", err
	}
	return result, nil
}

func (g *RedisGroup) LLen(ctx context.Context, key string) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LLen: key %s", key)
	count, err := node.client.LLen(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return count, nil
}

func (g *RedisGroup) LIndex(ctx context.Context, key string, index int64) (string, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return "", errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: LIndex: key %s index %d", key, index)
	val, err := node.client.LIndex(ctx, key, index).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return "", err
	}
	return val, nil
}

func (g *RedisGroup) SetNX(ctx context.Context, key string, val []byte, exp time.Duration) (bool, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return false, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: SetNX: key %s exp %v", key, exp)
	b, err := node.client.SetNX(ctx, key, val, exp).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return false, err
	}
	return b, nil
}

func (g *RedisGroup) ZScore(ctx context.Context, key string, member string) (float64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: ZScore: key %s member %s", key, member)
	score, err := node.client.ZScore(ctx, key, member).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return score, nil
}

func (g *RedisGroup) Ttl(ctx context.Context, key string) (time.Duration, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: TTL: key %s", key)
	ttl, err := node.client.TTL(ctx, key).Result()
	if err != nil {
		g.logger.Errorf(ctx, "err:%v", err)
		return 0, err
	}
	return ttl, nil
}

func (g *RedisGroup) SetBit(ctx context.Context, key string, offset int64, val int) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Set Bit: key %s offset %d val %d", key, offset, val)
	intCmd := node.client.SetBit(ctx, key, offset, val)
	if intCmd.Err() != nil {
		g.logger.Errorf(ctx, "err:%s", intCmd.Err())
		return 0, intCmd.Err()
	}

	return intCmd.Result()
}

func (g *RedisGroup) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: Set Bit: key %s offset %d", key, offset)
	intCmd := node.client.GetBit(ctx, key, offset)
	if intCmd.Err() != nil {
		g.logger.Errorf(ctx, "err:%s", intCmd.Err())
		return 0, intCmd.Err()
	}
	return intCmd.Result()
}

func (g *RedisGroup) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) (int64, error) {
	node := g.FindNodeForKey(key)
	if node == nil {
		return 0, errors.New("not found available redis node")
	}
	g.logger.Debugf(ctx, "redis: BitCount: key %s", key)
	intCmd := node.client.BitCount(ctx, key, bitCount)
	if intCmd.Err() != nil {
		g.logger.Errorf(ctx, "err:%s", intCmd.Err())
		return 0, intCmd.Err()
	}
	return intCmd.Result()
}

func (g *RedisGroup) FlushAll(ctx context.Context) (string, error) {
	g.logger.Debugf(ctx, "redis: FushAll")
	var intCmd *redis.StatusCmd
	for _, node := range g.nodes {
		intCmd = node.client.FlushAll(ctx)
		if intCmd.Err() != nil {
			g.logger.Errorf(ctx, "err:%s", intCmd.Err())
			return "", intCmd.Err()
		}
	}

	return intCmd.Result()
}