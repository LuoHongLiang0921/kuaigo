package redis

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/config"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"testing"
	"time"


	"github.com/stretchr/testify/assert"

)

func newRedis() config.IAdvanceCacheAdapter {
	ctx := context.Background()
	cfg := config.CacheConfig{
		Name:          "test",
		Mode:          StubMode,
		Type:          "redis",
		Addr:          "localhost:16379",
		DB:            0,
		AutoConnect:   false,
		PoolSize:      10,
		MaxRetries:    3,
		MinIdleConns:  100,
		DialTimeout:   ktime.Duration("1s"),
		ReadTimeout:   ktime.Duration("1s"),
		WriteTimeout:  ktime.Duration("1s"),
		IdleTimeout:   ktime.Duration("60s"),
		ReadOnly:      false,
		Debug:         false,
		EnableTrace:   false,
		SlowThreshold: ktime.Duration("250ms"),
		OnDialError:   "panic",
		Logger:        klog.KuaigoLogger,
	}
	redis := NewRedisAdapter(ctx, &cfg)
	return redis
}

func Test_redisAdapter_Decr(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	isBool := cfgAdapter.Decr("test.desr")
	assert.Equal(t, true, isBool, "test")
}
func Test_redisAdapter_Del(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	delNum := cfgAdapter.Del("test.desr")
	assert.NotZero(t, delNum, "delete test")

}

func Test_redisAdapter_DelWithErr(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	delNum, err := cfgAdapter.DelWithErr("test.del")
	assert.NoError(t, err, "del with err")
	assert.NotZero(t, delNum, "delete test")
}

func Test_redisAdapter_Exists(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	cfgAdapter.Set("test.exists", "2", 3*time.Second)
	is := cfgAdapter.Exists("test.exists")
	assert.Equal(t, true, is)
	time.Sleep(5 * time.Second)
	is = cfgAdapter.Exists("test.exists")
	assert.Equal(t, false, is)
}

func Test_redisAdapter_ExistsWithErr(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	cfgAdapter.Set("test.existswitherr", "2", 3*time.Second)
	is, err := cfgAdapter.ExistsWithErr("test.existswitherr")
	assert.Equal(t, true, is)
	assert.NoError(t, err, is)
	time.Sleep(5 * time.Second)
	is, err = cfgAdapter.ExistsWithErr("test.existswitherr")
	assert.Equal(t, false, is)
	assert.NoError(t, err, is)
}

func Test_redisAdapter_Expire(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	is := cfgAdapter.Set("test.expire", 2, 2*time.Second)
	assert.Equal(t, true, is, "set")
	is, err := cfgAdapter.Expire("test.expire", 5*time.Second)
	assert.NoError(t, err, "exipre")
	assert.Equal(t, true, is)
}

func TestRedisAdapter_HSet(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	is := cfgAdapter.HSet("test.hset", "s", 1)
	assert.Equal(t, true, is, "set")
	is, err := cfgAdapter.Expire("test.hset", 5*time.Second)
	assert.NoError(t, err, "exipre")
	assert.Equal(t, true, is)
}
func TestRedisAdapter_HMSet(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	is := cfgAdapter.HMSet("test.hset", map[string]interface{}{"t": "t", "t1": 1}, 2*time.Second)
	assert.Equal(t, true, is)
}
func TestRedisAdapter_HGet(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	is := cfgAdapter.HMSet("test.hget", map[string]interface{}{"t": "t", "t1": 1}, 2*time.Second)
	assert.Equal(t, true, is, "hget")
	r, err := cfgAdapter.HGet("test.hget", "t")
	assert.NoError(t, err)
	assert.Equal(t, "t", r)
}
func TestRedisAdapter_MGets(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	cfgAdapter.Set("t.mget.1", 1, 3*time.Second)
	cfgAdapter.Set("t.mget.2", 2, 3*time.Second)
	cfgAdapter.Set("t.mget.3", 3, 3*time.Second)
	r, err := cfgAdapter.MGet("t.mget.1", "t.mget.2", "t.mget.3")
	assert.NoError(t, err, "")
	assert.Equal(t, []string{"1", "2", "3"}, r)
}

func TestRedisAdapter_SetNx(t *testing.T) {
	cfgAdapter := newRedis()
	defer cfgAdapter.Close()
	is := cfgAdapter.Set("test.nx", 5, 2*time.Second)
	assert.Equal(t, true, is)
	is = cfgAdapter.SetNx("test.nx", 6, 2*time.Second)
	is = assert.Equal(t, false, is, "set nx 2")
	time.Sleep(4 * time.Second)
	assert.Equal(t, true, is, "set nx 3")
}
