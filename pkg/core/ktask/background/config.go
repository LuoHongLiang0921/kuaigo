package background

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

// Config 配置
type Config struct {
	// Name 任务名字
	Name string
}

// RawConfig
// 	@Description
//	@Param ctx
//	@Param key
// 	@Return Config
func RawConfig(ctx context.Context, key string) Config {
	var config Config
	if err := conf.UnmarshalKey(key, &config); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panicf("key %v unmarshal err:%v", key, err)
	}
	return config
}

// Build
// 	@Description
// 	@Receiver c Config
// 	@Return *Background
func (c Config) Build() *Background {
	return &Background{
		config: &c,
	}
}
