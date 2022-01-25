package config

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"
)

const (
	cacheConfigPrefix = "caches."
	//ClusterMode using clusterClient
	ClusterMode string = "cluster"
	//StubMode using redisClient
	StubMode       string = "stub"
	RootDefaultKey        = "caches"
)

// CacheConfig for redis, contains RedisStubConfig and RedisClusterConfig
type CacheConfig struct {
	//名称
	Name string `json:"name" yaml:"name"`
	//类型
	Type string `json:"type" yaml:"type"`
	// Addrs 实例配置地址
	Addrs []string `json:"addrs" yaml:"addrs"`
	// Addr stubConfig 实例配置地址
	Addr string `json:"addr" yaml:"addr"`
	// Mode Redis模式 cluster|stub
	Mode string `json:"mode" yaml:"mode"`
	// Password 密码
	Password string `json:"password" yaml:"password"`
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int `json:"db" yaml:"db"`
	//启动是否自动连接
	AutoConnect bool `json:"autoConnect" yaml:"autoConnect"`
	// PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	PoolSize int `json:"poolSize" yaml:"poolSize" `
	// MaxRetries 网络相关的错误最大重试次数 默认8次
	MaxRetries int `json:"maxRetries" yaml:"maxRetries"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `json:"minIdleConns" yaml:"minIdleConns"`
	// DialTimeout 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	// ReadTimeout 读超时 默认3s
	ReadTimeout time.Duration `json:"readTimeout" yaml:"readTimeout"`
	// WriteTimeout 读超时 默认3s
	WriteTimeout time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	// IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	IdleTimeout time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
	// Debug开关
	Debug bool `json:"debug" yaml:"debug"`
	// ReadOnly 集群模式 在从属节点上启用读模式
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
	// 是否开启链路追踪，开启以后。使用DoContext的请求会被trace
	EnableTrace bool `json:"enableTrace" yaml:"enableTrace"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
	// OnDialError panic|error
	OnDialError string `json:"level" yaml:"level"`
	Logger      *klog.Logger

	latestDsn string
	change    chan struct{}
}

// IsConfigChange
// 	@Description
// 	@Receiver c
// 	@Return <-chan
func (c *CacheConfig) IsConfigChange() <-chan struct{} {
	return c.change
}

func (c *CacheConfig) setOnChange(key string) {
	configKey := c.getConfigKey(key)
	conf.OnChange(func(cfg *conf.Configuration) {
		addrRoot := configKey + ".addr"
		dsnStr := cfg.GetString(addrRoot)
		klog.Debugf("%s change, result %ss", addrRoot, dsnStr)
		if dsnStr != "" && dsnStr != c.latestDsn {
			c.change <- struct{}{}
		}
		c.latestDsn = dsnStr
	})
}

func (c *CacheConfig) getConfigKey(key string) string {
	return cacheConfigPrefix + key
}

// GetConfig
// 	@Description 获取默认配置
//  @Param ctx 上下文Context
// 	@Return Config 设置默认值后的配置
func GetConfig(ctx context.Context, key string) *CacheConfig {
	config := &CacheConfig{
		Name:          key,
		Type:          "redis",
		DB:            0,
		AutoConnect:   false,
		PoolSize:      10,
		MaxRetries:    3,
		MinIdleConns:  100,
		DialTimeout:   ktime.Duration("1s"),
		ReadTimeout:   ktime.Duration("1s"),
		WriteTimeout:  ktime.Duration("1s"),
		IdleTimeout:   ktime.Duration("60s"),
		ReadOnly:      false,
		Debug:         false,
		EnableTrace:   false,
		SlowThreshold: ktime.Duration("250ms"),
		OnDialError:   "panic",
		Logger:        klog.KuaigoLogger,
		change:        make(chan struct{}, 1),
	}
	config.latestDsn = config.Addr
	configKey := config.getConfigKey(key)
	if err := conf.UnmarshalKey(configKey, config); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panic("unmarshal redisConfig",
			klog.String("key", configKey),
			klog.Any("redisConfig", config),
			klog.String("error", err.Error()))
	}
	config.latestDsn = config.Addr
	config.setOnChange(key)
	return config
}
