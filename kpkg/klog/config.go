package klog

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"time"

	"github.com/go-redis/redis"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int
	MaxAge    int
	MaxBackup int
	// 日志磁盘刷盘间隔
	Interval      time.Duration
	CallerSkip    int
	Async         bool
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	RedisCore     zapcore.Core
	Debug         bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
	RedisURL      string
	RedisPass     string
	LogKey        string
	redisClient   *redis.Client
	//是否放到支持redis协议的es中
	Save bool
	// 业务日志类型
	//
	LoggerType string
	// yw
	ServiceName string
}

// Filename ...
func (config *Config) Filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		panic(err)
	}
	config.configKey = key
	return config
}

// StdConfig tabby Standard logger config
func StdConfig(name string) *Config {
	return RawConfig("tabby.logger." + name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:          "bizLogger.log",
		Dir:           ".",
		Level:         "info",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    1,
		AddCaller:     false,
		Async:         true,
		Queue:         false,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: DefaultZapConfig(),
	}
}

func (config Config) getContext() context.Context {
	return context.TODO()
}

// Build ...
func (config Config) Build() *Logger {
	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapConfig()
	}
	if config.Debug {
		config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	if config.redisClient == nil && config.RedisURL != "" && config.Save {
		client := redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     config.RedisURL,
			Password: config.RedisPass,
		})
		config.redisClient = client
	}
	logger := newLogger(&config)
	if config.configKey != "" {
		logger.AutoLevel(config.configKey + ".level")
	}
	return logger
}
