package rabbitmq

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kpool"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"

	"github.com/streadway/amqp"
)

type Config struct {
	Address   string
	MaxIdle   int
	MaxActive int
	// IdleTimeout 池中空闲时间，每次获取都会清理池中的连接，每次都会从
	IdleTimeout time.Duration
	// MaxConnLifetime 最大存活时间，
	MaxConnLifetime time.Duration
	// Mode 队列模式,kafka,rocketmq,rabbitmq
	Mode string
	// RunType
	RunType mq.RunType
	ConsumeConfig
	PublishConfig
}

type PublishConfig struct {
	Exchange string
}

type ConsumeConfig struct {
	// 消费者需要 队列名字
	Queue       string
	ConsumerTag string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	Args        map[string]interface{}

	IsNackRequeue bool
}

// RawConfig
// 	@Description 载入配置
//	@Param key mq 配置文件中key
// 	@Return Config 配置结构体
func RawConfig(key string) Config {
	config := getDefaultConfig()
	if err:= conf.UnmarshalKey(key, &config); err != nil {
		klog.Panic("unmarshal MQConfig",
			klog.String("key", key),
			klog.Any("mqConfig", config),
			klog.String("error", err.Error()))
	}
	return config
}

func getDefaultConfig() Config {
	return Config{
		Mode: constant.ModeRabbitmq,
	}
}

func (c *Config) setDefaultConfig() {
	if c.MaxIdle == 0 {
		c.MaxIdle = 100
	}
	if c.MaxActive == 0 {
		c.MaxActive= 100
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = ktime.Duration("1m")
	}
	if c.MaxConnLifetime == 0 {
		c.MaxConnLifetime = ktime.Duration("2m")
	}
}

// Build
// 	@Description 构建rabbitmq 实例池
// 	@Receiver c Config
// 	@Return *Client rabbitmq 实例池
// 	@Return error 错误
func (c Config) Build(ctx context.Context) *Client {
	c.mustValidConfig()
	c.setDefaultConfig()
	p := &Client{
		cfg:     &c,
		closing: make(chan struct{}),
		logger:  klog.KuaigoLogger,
	}
	pool := &kpool.Pool{
		Name: "rabbitmq",
		Dial: func() (kpool.Conn, error) {
			aqConn, err := amqp.Dial(p.cfg.Address)
			if err != nil {
				klog.WithContext(ctx).Errorf("Dial addr:%s failed, err: %v", p.cfg.Address, err)
				return nil, err
			}
			return aqConn, nil
		},
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		IdleTimeout:     c.IdleTimeout,
		MaxConnLifetime: c.MaxConnLifetime,
	}
	p.pool = pool
	return p
}

func (c Config) mustValidConfig() {
	if c.Mode == "" {
		klog.KuaigoLogger.Panic("config mode must not empty")
	}
	if c.RunType == "" {
		klog.KuaigoLogger.Panic("config run type must not empty")
	}
	if len(c.Address) <= 0 {
		klog.KuaigoLogger.Panic("config address must not empty")
	}
}
