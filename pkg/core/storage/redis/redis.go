// @Description

package redis

import "github.com/go-redis/redis"

//TODO 引入redis统一错误码

//Redis client (cmdable and config)
type Redis struct {
	Config *Config
	Client redis.Cmdable
}

// Cluster try to get a redis.ClusterClient
func (r *Redis) Cluster() *redis.ClusterClient {
	if c, ok := r.Client.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

//Stub try to get a redis.Client
func (r *Redis) Stub() *redis.Client {
	if c, ok := r.Client.(*redis.Client); ok {
		return c
	}
	return nil
}

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
