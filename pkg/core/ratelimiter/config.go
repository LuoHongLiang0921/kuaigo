// @Description 限流控制器配置

package ratelimiter

import (
	kconf "github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

// Config 配置信息
type Config struct {
	// RateLimiter 多个限流策略
	Rule map[string]*RLConfig
	// RateLimiterPaths
	// 不同限流策略下多个path
	// 同一path不能在多个限流策略中,如在不同策略中,会覆盖只保留一个！
	Path  map[string][]string
	redis string
	// resourceAns
	resourceAns map[string][]string
}

func (c *Config) setRule() {
	resourceAns := make(map[string][]string, 0)
	for an, ls := range c.Path {
		for _, l := range ls {
			resourceAns[l] = append(resourceAns[l], an)
		}
	}
	c.resourceAns = resourceAns
}

// RawConfig
//  @Description: 获取rate配置信息
//  @Param key
//  @Return *Config
func RawConfig(key string) *Config {
	var appCfg Config
	err := kconf.UnmarshalKey(key, &appCfg)
	if err != nil {
		klog.Panic("unmarshal RateLimiter config")
	}
	appCfg.setRule()
	// 加上次方法可以使apollo实时生效
	kconf.OnChange(func(cfg *kconf.Configuration) {
		err := kconf.UnmarshalKey(key, &appCfg)
		if err != nil {
			klog.Panic("unmarshal RateLimiter config")
		}
		appCfg.setRule()
	})
	return &appCfg
}

// Build
//  @Description:构建限流工具
//  @Receiver c
//  @Return *RateLimiter
func (c *Config) Build() *RateLimiter {
	return &RateLimiter{
		cfg: c,
	}
}
