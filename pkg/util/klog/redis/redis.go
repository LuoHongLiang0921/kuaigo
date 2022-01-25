// @Description

package redis

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"sync"
)

const (
	OutputRedis = "redis"
)

var once sync.Once

// RegisterOutputCreatorHandler
// 	@Description
func RegisterOutputCreatorHandler() {
	once.Do(func() {
		manager.Register(OutputRedis, func(cfg interface{}) []zap.Core {
			if redisCfg, ok := cfg.(Config); ok {
				redisCore := redisCfg.Build()
				errRedisCore := redisCfg.BuildAlter("error", nil)
				zapCores := make([]zap.Core, 0, 2)
				if redisCore != nil {
					zapCores = append(zapCores, redisCore)
				}
				if errRedisCore != nil {
					zapCores = append(zapCores, redisCore)
				}
				return zapCores
			}
			return nil
		})
	})
}
