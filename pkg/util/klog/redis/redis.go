// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/6/30 9:43 上午

package redis

import (
	"sync"

	"git.bbobo.com/framework/tabby/pkg/util/xlog/manager"
	"git.bbobo.com/framework/tabby/pkg/util/xlog/zap"
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
