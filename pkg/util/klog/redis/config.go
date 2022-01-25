// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/6/29 9:43 上午

package redis

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.bbobo.com/framework/tabby/pkg/util/xlog/buffer"

	"git.bbobo.com/framework/tabby/pkg/util/xtime"

	"git.bbobo.com/framework/tabby/pkg/conf"
	"git.bbobo.com/framework/tabby/pkg/defers"
	"git.bbobo.com/framework/tabby/pkg/util/xlog/zap"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	// 	Addr string redis 连接串
	Addr string
	// Password 密码
	Password string

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Key Redis 协议总key
	Key string
	// Level 日志等级
	Level string
	// redis 来源配置
	Source string

	//Async 异步
	Async bool `yaml:"async"`
	// FlushInterval
	FlushInterval string `yaml:"flushInterval"`
	// BufferSize 缓冲大小
	BufferSize int `yaml:"bufferSize"`

	// redisWriter redisWrite 写入组件
	redisWriter *RedisLog `mapstructure:"-"`
	once        sync.Once

	parentKey string           `mapstructure:"-"`
	lv        *zap.AtomicLevel `mapstructure:"-"`
}

// LoaSourceConfig
// 	@Description 载入 source 字段的配置
// 	@Receiver c
// 	@Return error
func (c *Config) LoaSourceConfig() error {
	if c.Source == "" {
		return errors.New("source is empty")
	}
	connectInfo := getDefaultRedisConfig()
	redisConn := conf.Get(c.Source)
	decodeCofig := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &connectInfo,
	}
	decoder, err := mapstructure.NewDecoder(&decodeCofig)
	if err != nil {
		return err
	}
	err = decoder.Decode(redisConn)
	if err != nil {
		return err
	}
	c.Addr = connectInfo.Addr
	c.Password = connectInfo.Password
	c.WriteTimeout = connectInfo.WriteTimeout
	c.ReadTimeout = connectInfo.ReadTimeout
	c.DialTimeout = connectInfo.DialTimeout
	c.IdleTimeout = connectInfo.IdleTimeout
	return nil
}

func getDefaultRedisConfig() Config {
	return Config{
		DialTimeout:  xtime.Duration("1s"),
		ReadTimeout:  xtime.Duration("1s"),
		WriteTimeout: xtime.Duration("1s"),
		IdleTimeout:  xtime.Duration("60s"),
	}
}

func (c *Config) SetDefaultConfig() {
	if c.Async {
		if c.FlushInterval == "" {
			c.FlushInterval = "5s"
		}
		if c.BufferSize == 0 {
			c.BufferSize = 256 * 1024
		}
	}

}

// SetRedisLog
// 	@Description 设置 redis 组件
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetRedisLog() *Config {
	c.once.Do(func() {
		if c.Addr != "" {
			client := redis.NewClient(&redis.Options{
				Network:      "tcp",
				Addr:         c.Addr,
				Password:     c.Password,
				WriteTimeout: c.WriteTimeout,
				ReadTimeout:  c.ReadTimeout,
				IdleTimeout:  c.IdleTimeout,
				DialTimeout:  c.DialTimeout,
			})
			redisLog := NewRedisLog(c.Key, client)
			c.redisWriter = redisLog
		}
	})
	return c
}

// SetParent
// 	@Description 设置父级key
// 	@Receiver c
//	@Param k
// 	@Return *Config
func (c *Config) SetParent(k string) *Config {
	c.parentKey = k
	return c
}

// SetAutoLevel
// 	@Description 设置auto level
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetAutoLevel() *Config {
	var lv zap.Level

	err := lv.Set(c.Level)
	if err != nil {
		panic(err)
	}
	alv := zap.NewAtomicLevelAt(lv)
	err = alv.UnmarshalText([]byte(c.Level))
	if err != nil {
		panic(err)
	}
	c.lv = &alv
	if c.parentKey == "" {
		return c
	}
	conf.OnChange(func(config *conf.Configuration) {
		lvText := strings.ToLower(config.GetString(c.parentKey + ".output.redis.level"))
		if lvText != "" {
			err := c.lv.UnmarshalText([]byte(lvText))
			if err != nil {
				return
			}
		}
	})
	return c
}

// BuildAlter
// 	@Description 构建错误报警日志实例，只有错误日志级别才报警
// 	@Receiver c
//	@Param alterKey 报警key
// 	@Return zap.Core
func (c *Config) BuildAlter(alterKey string, enablerFunc zap.LevelEnablerFunc) zap.Core {
	c.SetRedisLog()
	redisLog := c.redisWriter
	if redisLog == nil {
		return nil
	}
	encoder := func() zap.Encoder {
		return zap.NewJSONEncoder(*getRedisZapConfig())
	}

	errRedisLog := redisLog.SetLogKey(alterKey)
	errWs := zap.AddSync(errRedisLog)
	if c.Async {
		var close buffer.CloseFunc
		errWs, close = buffer.Buffer(errWs, c.BufferSize, xtime.Duration(c.FlushInterval))
		defers.Register(close)
	}
	errorCore := zap.NewCore(encoder(), errWs, zap.LevelEnablerFunc(func(level zap.Level) bool {
		return level == zap.ErrorLevel
	}))
	return errorCore
}

// Build
// 	@Description 构建日志实例
// 	@Receiver c
// 	@Return zap.Core
func (c *Config) Build() zap.Core {
	c.SetRedisLog()
	redisLog := c.redisWriter
	if redisLog == nil {
		return nil
	}
	rws := zap.AddSync(redisLog)
	encoder := func() zap.Encoder {
		return zap.NewJSONEncoder(*getRedisZapConfig())
	}
	if c.lv == nil {
		panic(fmt.Errorf("%s atom level is empty", c.parentKey))
	}
	if c.Async {
		var close buffer.CloseFunc
		rws, close = buffer.Buffer(rws, c.BufferSize, xtime.Duration(c.FlushInterval))
		defers.Register(close)
	}
	redisZapCore := zap.NewCore(encoder(), rws, zap.LevelEnablerFunc(func(level zap.Level) bool {
		if level == zap.ErrorLevel {
			return false
		}
		return c.lv.Enabled(level)
	}))
	defers.Register(redisLog.Close)
	return redisZapCore
}

func getRedisZapConfig() *zap.EncoderConfig {
	return &zap.EncoderConfig{
		TimeKey:          "timestamp",
		LevelKey:         zap.OmitKey,
		NameKey:          zap.OmitKey,
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stack",
		LineEnding:       zap.DefaultLineEnding,
		EncodeLevel:      zap.LowercaseLevelEncoder,
		EncodeTime:       zap.EpochMillisTimeEncoder,
		EncodeDuration:   zap.SecondsDurationEncoder,
		EncodeCaller:     zap.ShortCallerEncoder,
		EncodeName:       zap.FullNameEncoder,
		ConsoleSeparator: " ",
	}
}
