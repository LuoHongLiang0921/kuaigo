package config

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type (
	Z                  = redis.Z
	Pipeline           = redis.Pipeline
	StringCmd          = redis.StringCmd
	StringSliceCmd     = redis.StringSliceCmd
	IntCmd             = redis.IntCmd
	DurationCmd        = redis.DurationCmd
	TimeCmd            = redis.TimeCmd
	BoolCmd            = redis.BoolCmd
	BoolSliceCmd       = redis.BoolSliceCmd
	FloatCmd           = redis.FloatCmd
	StringStringMapCmd = redis.StringStringMapCmd
	StringIntMapCmd    = redis.StringIntMapCmd
	StringStructMapCmd = redis.StringStructMapCmd
	GeoRadiusQuery     = redis.GeoRadiusQuery
	GeoLocation        = redis.GeoLocation
	ZRangeBy           = redis.ZRangeBy
)

var (
	Nil = redis.Nil
)

// ICache 普通操作实现接口
type ICache interface {
	WithContext(ctx context.Context) ICache
	Get(key string) string
	GetRaw(key string) ([]byte, error)
	MGet(keys ...string) ([]string, error)
	MGets(keys []string) ([]interface{}, error)
	Set(key string, value interface{}, expire time.Duration) bool
	SetWithErr(key string, value interface{}, expire time.Duration) error
	SetNx(key string, value interface{}, expiration time.Duration) bool
	SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error)
	Incr(key string) bool
	IncrWithErr(key string) (int64, error)
	IncrBy(key string, increment int64) (int64, error)
	Decr(key string) bool
	Del(key ...string) int64
	DelWithErr(key string) (int64, error)
	Exists(key string) bool
	ExistsWithErr(key string) (bool, error)
	Expire(key string, expiration time.Duration) (bool, error)
	TTL(key string) (int64, error)
}

// IAdvanceCache 高级操作实现接口
type IAdvanceCache interface {
	ICache
	WithAdvanceContext(ctx context.Context) IAdvanceCache
	HGetAll(key string) map[string]string
	HGet(key string, fields string) (string, error)
	HMGet(key string, fields []string) []string
	HMGetMap(key string, fields []string) map[string]string
	HMSet(key string, hash map[string]interface{}, expire time.Duration) bool
	HSet(key string, field string, value interface{}) bool
	HDel(key string, field ...string) bool
	Scan(cursor uint64, match string, count int64) ([]string, error)
	Type(key string) (string, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]Z, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRevRank(key string, member string) (int64, error)
	ZRevRangeByScore(key string, opt ZRangeBy) ([]string, error)
	ZRevRangeByScoreWithScores(key string, opt ZRangeBy) ([]Z, error)
	ZCard(key string) (int64, error)
	ZScore(key string, member string) (float64, error)
	ZAdd(key string, members ...Z) (int64, error)
	ZCount(key string, min, max string) (int64, error)
	HIncrBy(key string, field string, incr int) int64
	HIncrByWithErr(key string, field string, incr int) (int64, error)
	LPush(key string, values ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	RPop(key string) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	LLen(key string) int64
	LLenWithErr(key string) (int64, error)
	LRem(key string, count int64, value interface{}) int64
	LIndex(key string, idx int64) (string, error)
	LTrim(key string, start, stop int64) (string, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRemRangeByScore(key string, min, max string) (int64, error)
	ZRem(key string, members ...interface{}) (int64, error)
	SAdd(key string, member ...interface{}) (int64, error)
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	HKeys(key string) []string
	HLen(key string) int64
	GeoAdd(key string, location *GeoLocation) (int64, error)
	GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) ([]GeoLocation, error)
}

// ICacheAdapter 缓存适配器接口
type ICacheAdapter interface {
	ICache
	Open(ctx context.Context) ICache
	GetClient() ICache
	Close() (err error)
}

// IAdvanceCacheAdapter 缓存高级适配器接口
type IAdvanceCacheAdapter interface {
	IAdvanceCache
	Open(ctx context.Context) ICache
	GetClient() ICache
	Close() (err error)
}
