package cache

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/adapter/redis"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/adapter/rediscluster"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/config"
	"sync"
)

var cacheManagerInstance *cacheManager
var instanceOnce sync.Once
var sum sync.RWMutex

type cacheManager struct {
	caches map[string]config.ICache
}

// GetCacheManagerInstance
//  @Description 获取全局单例，如果new多个，永远指向最后一次Build的实例
func GetCacheManagerInstance() *cacheManager {
	instanceOnce.Do(func() {
		cacheManagerInstance = new(cacheManager)
		cacheManagerInstance.caches = make(map[string]config.ICache, 0)
	})
	return cacheManagerInstance
}

// GetCache
// 	@Description 获取缓存操作类
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param conf 数据库配置key
// 	@Return IDataBaseAdapter
func (df *cacheManager) GetCache(ctx context.Context, conf string) config.ICache {
	if db, ok := df.caches[conf]; ok {
		return db
	}
	return df.buildCache(ctx, conf)
}

// GetAdvanceCache
// 	@Description 获取高级缓存操作类
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param conf 数据库配置key
// 	@Return IDataBaseAdapter
func (df *cacheManager) GetAdvanceCache(ctx context.Context, conf string) config.IAdvanceCache {
	if db, ok := df.caches[conf]; ok {
		return db.(config.IAdvanceCache)
	}
	return df.buildCache(ctx, conf).(config.IAdvanceCache)
}

// buildCache
// 	@Description 构造缓存内部方法
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param config 数据库配置key
// 	@Return IDatabase
func (df *cacheManager) buildCache(ctx context.Context, conf string) config.ICache {
	sum.Lock()
	if _, ok := df.caches[conf]; !ok {
		df.caches[conf] = df.buildCacheAdapter(ctx, conf)
	}
	sum.Unlock()
	return df.caches[conf]
}

// buildCacheAdapter
// 	@Description 构造缓存适配器内部方法
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param config 数据库配置key
// 	@Return IDataBaseAdapter
func (df *cacheManager) buildCacheAdapter(ctx context.Context, conf string) config.ICacheAdapter {
	dbConfig := config.GetConfig(ctx, conf)
	if dbConfig.Type == "redis" {
		return redis.NewRedisAdapter(ctx, dbConfig)
	} else if dbConfig.Type == "redisCluster" {
		return rediscluster.NewRedisClusterAdapter(ctx, dbConfig).(config.ICacheAdapter)
	}
	panic("tabby not support " + dbConfig.Type + " cache")
}
