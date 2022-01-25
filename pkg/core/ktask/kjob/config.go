package kjob

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

type Config struct {
	Name string
}

// RawConfig
// 	@Description 获取配置
//	@Param ctx 上下文
//	@Param key 一次性任务对应的key
// 	@Return Config 配置实例
func RawConfig(ctx context.Context, key string) Config {
	var config Config
	if err := conf.UnmarshalKey(key, &config); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panicf("key %v unmarshal err:%v", key, err)
	}
	return config
}

// Build
// 	@Description 实例化一次性任务
// 	@Receiver c Config
//	@Param ctx 上下文
// 	@Return *XJob 一次性job 实例
func (c Config) Build(ctx context.Context) *XJob {
	return &XJob{
		config: &c,
		ctx:    ctx,
		closed: make(chan struct{}),
	}
}
