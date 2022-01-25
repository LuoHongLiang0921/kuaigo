package rocketmq

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

type Config struct {
	// Mode 队列模式,kafka,rocketmq,rabbitmq
	Mode string
	// RunType
	RunType   mq.RunType
	Brokers   []string
	Topic     string
	GroupName string
	// Retry 重试次数
	Retry int
	// Async 是否异步
	Async bool
}

// RawConfig
// 	@Description 实例化配置
//	@Param ctx 上下文
//	@Param key 配置key
// 	@Return *Config 实例后的配置
func RawConfig(ctx context.Context, key string) *Config {
	cfg := getDefaultConfig()
	if err := conf.UnmarshalKey(key, &cfg); err != nil {
		klog.WithContext(ctx).Panicf("kafka config err:%v", err)
	}
	return &cfg
}

func getDefaultConfig() Config {
	return Config{
		Mode:  constant.ModeRocketmq,
		Async: true,
	}
}

// Build
// 	@Description 实例rocket mq 实例
// 	@Receiver c Config
//	@Param ctx 上下文
// 	@Return *Client
func (c Config) Build(ctx context.Context) *Client {
	client := NewClient(ctx, &c)
	return client
}
