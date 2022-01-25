// @Description

package redis

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"


	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
)

const (
	//ClusterMode using cluster client
	ClusterMode = "cluster"
	//StubMode using redis client
	StubMode       = "stub"
	RootDefaultKey = "caches"
)

// Config for redis, contains RedisStubConfig and RedisClusterConfig
type Config struct {
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
	// 是否开启链路追踪，开启以后。使用DoCotext的请求会被trace
	EnableTrace bool `json:"enableTrace" yaml:"enableTrace"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
	// OnDialError panic|error
	OnDialError string `json:"level" yaml:"level"`
	logger      *klog.Logger
}

// DefaultRedisConfig
// 	@Description 获取默认配置
// 	@Return Config 设置默认值后的配置
func DefaultRedisConfig() Config {
	return Config{
		DB:            0,
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
		logger:        klog.KuaigoLogger,
	}
}

// StdRedisConfig
// 	@Description: 取为 "caches.{name}" 的配置，其中为 name 为某个实例别名
//	@Param name 配置别名
// 	@Return Config 载入后的配置
func StdRedisConfig(name string) Config {
	return RawRedisConfig(RootDefaultKey + name)
}

// RawRedisConfig
// 	@Description: 载入配置
//	@Param key 全名称key
// 	@Return Config 载入后的配置
func RawRedisConfig(key string) Config {
	var config = DefaultRedisConfig()

	if err := conf.UnmarshalKey(key, &config); err != nil {
		klog.KuaigoLogger.Panic("unmarshal redisConfig",
			klog.String("key", key),
			klog.Any("redisConfig", config),
			klog.String("error", err.Error()))
	}
	return config
}

// MultiCache
// 	@Description: 载入并实例化多数据源缓存
//	@Param ctx 上下文
//	@Param key 根 key 名字
// 	@Return map[string]*Redis 多数据源缓存
func MultiCache(ctx context.Context, key string) map[string]*Redis {
	dataV := conf.GetStringMap(key)
	m := make(map[string]*Config)
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &m,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("redis"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	err = decoder.Decode(dataV)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("redis"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	dbMap := make(map[string]*Redis)
	for k, v := range m {
		db := v.Build()
		dbMap[k] = db
	}
	return dbMap
}

// Build
// 	@Description: 初始化redis 客户端
// 	@receiver config 配置
// 	@Return *Redis 实例化后的redis 实例
func (config Config) Build() *Redis {
	config.mustValidConfig()
	var client redis.Cmdable
	switch config.Mode {
	case ClusterMode:
		if len(config.Addrs) == 1 {
			config.logger.Warn("redis config has only 1 address but with cluster mode")
		}
		client = config.buildCluster()
	case StubMode:
		client = config.buildStub()
	default:
		config.logger.Panic("redis mode must be one of (stub, cluster)")
	}
	return &Redis{
		Config: &config,
		Client: client,
	}
}

func (config Config) mustValidConfig() {
	if config.Mode == "" {
		config.logger.Panic("redis mode must be one of (stub, cluster),but is empty")
	}
}

func (config Config) buildStub() *redis.Client {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})

	if err := stubClient.Ping().Err(); err != nil {
		switch config.OnDialError {
		case "panic":
			config.logger.Panic("dial redis fail", klog.Any("err", err), klog.Any("config", config))
		default:
			config.logger.Error("dial redis fail", klog.Any("err", err), klog.Any("config", config))
		}
	}

	return stubClient

}

func (config Config) buildCluster() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		MaxRedirects: config.MaxRetries,
		ReadOnly:     config.ReadOnly,
		Password:     config.Password,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	if err := clusterClient.Ping().Err(); err != nil {
		switch config.OnDialError {
		case "panic":
			config.logger.Panic("start cluster redis", klog.Any("err", err))
		default:
			config.logger.Error("start cluster redis", klog.Any("err", err))
		}
	}
	return clusterClient
}
