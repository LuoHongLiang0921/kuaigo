package model

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/config"
)

type BaseCacheModel struct {
	BaseModel
	cacheConfig  string
	Cache        config.ICache
	AdvanceCache config.IAdvanceCache
}

// BuildCache
//  @Description 构造函数，New完之后必须调用。初始化前置资源
func (m *BaseCacheModel) BuildCache(ctx context.Context) *BaseCacheModel {
	m.Build(ctx)
	m.Cache = cache.GetCacheManagerInstance().GetCache(ctx, m.cacheConfig)
	m.AdvanceCache = cache.GetCacheManagerInstance().GetAdvanceCache(ctx, m.cacheConfig)
	return m
}

// WithCacheConfig
// 	@Description 缓存配置类型
// 	@Receiver BaseModel
//	@Param config 缓存配置key
// 	@Return *BaseModel
func (m *BaseCacheModel) WithCacheConfig(configCache string) *BaseCacheModel {
	m.cacheConfig = configCache
	return m
}
