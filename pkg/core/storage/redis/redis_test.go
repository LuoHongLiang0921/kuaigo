// @Description

package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	// TODO(gorexlv): add redis ci
	redisConfig := DefaultRedisConfig()
	redisConfig.Addrs = []string{"localhost:6379"}
	redisConfig.Mode = StubMode
	redisClient := redisConfig.Build()
	err := redisClient.Client.Ping().Err()
	if err != nil {
		t.Errorf("redis ping failed:%v", err)
	}
	st := redisClient.Stub().PoolStats()
	t.Logf("running status %+v", st)
	err = redisClient.Close()
	if err != nil {
		t.Errorf("redis close failed:%v", err)
	}
	st = redisClient.Stub().PoolStats()
	t.Logf("close status %+v", st)
}
