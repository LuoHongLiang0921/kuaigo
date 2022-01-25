// @Description 日志配置
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package klog

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/console"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/file"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/redis"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"github.com/mitchellh/mapstructure"
)

// Config 日志配置项
type Config struct {
	// ConfigVersion 配置版本
	ConfigVersion string
	// Deprecated: 输出类型 file:1, redis:2,http:4
	store int
	// Deprecated: Store 输出类型 file redis http
	Store []string
	// ServiceName 服务名
	ServiceName string

	//Async 异步
	Async bool
	//Level 日志初始等级
	Level string
	//AddCaller 是否添加调用者信息
	AddCaller bool
	// CallerSkip
	CallerSkip int

	// LoggerType 业务日志类型
	LoggerType string
	//Output 日志输出源配置
	Output map[string]interface{}

	//Alias 配置别名 full
	Alias string `yaml:"-" mapstructure:"-"`
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config Config
	if err := conf.UnmarshalKey(key, &config); err != nil {
		panic(err)
	}
	config.Alias = key
	return &config
}

// StdConfig
// 	@Description: 标准输入，省略 "logging" 前缀
//	@Param name logging 中日志别名
// 	@return *Config
func StdConfig(name string) *Config {
	return RawConfig("logging." + name)
}

// WithConfigVersion
// 	@Description 设置配置版本
// 	@Receiver c Config
//	@Param v 版本
// 	@Return *Config
func (c *Config) WithConfigVersion(v string) *Config {
	c.ConfigVersion = v
	return c
}

// WithServiceName
// 	@Description 设置服务名
// 	@Receiver c
//	@Param v
// 	@Return *Config
func (c *Config) WithServiceName(v string) *Config {
	c.ServiceName = v
	return c
}

// WithLoggerType
// 	@Description 设置logger typer
// 	@Receiver c
//	@Param t
// 	@Return *Config
func (c *Config) WithLoggerType(t string) *Config {
	c.LoggerType = t
	return c
}

// RegisterOutput
// 	@Description 注册输出源
// 	@Receiver c
// 	@Return Config
func (c Config) RegisterOutput() Config {
	console.RegisterOutputCreatorHandler()
	file.RegisterOutputCreatorHandler()
	redis.RegisterOutputCreatorHandler()
	return c
}

func getDefaultConfig(parentKey string) *Config {
	var cfg Config
	var consoleCfg console.Config
	consoleCfg.SetDefaultConfig()
	cfg.Output = map[string]interface{}{
		console.OutputConsole: consoleCfg,
	}
	cfg.Alias = parentKey
	return &cfg
}

// Build
// 	@Description  根据配置实例日志对象
// 	@receiver c Config
// 	@return *Logger
func (c Config) Build() *Logger {
	var cores []zap.Core
	var zapOptions []zap.Option
	for outputType, v := range c.Output {
		var outputCfg interface{}
		switch outputType {
		case console.OutputConsole:
			var consoleCfg console.Config
			err := mapstructure.Decode(v, &consoleCfg)
			if err != nil {
				panic(err)
			}
			consoleCfg.SetDefaultConfig()
			if consoleCfg.Level == "" {
				consoleCfg.Level = c.Level
			}
			consoleCfg.SetParent(c.Alias).SetAutoLevel()
			outputCfg = consoleCfg
		case file.OutputFile:
			var fileCfg file.Config
			err := mapstructure.Decode(v, &fileCfg)
			if err != nil {
				panic(err)
			}
			fileCfg.SetDefaultConfig()
			if fileCfg.Level == "" {
				fileCfg.Level = c.Level
			}
			fileCfg.SetParent(c.Alias).SetAutoLevel()
			outputCfg = fileCfg
		case redis.OutputRedis:
			var redisCfg redis.Config
			err := mapstructure.Decode(v, &redisCfg)
			if err != nil {
				panic(err)
			}
			err = redisCfg.LoaSourceConfig()
			if err != nil {
				panic(err)
			}
			redisCfg.SetDefaultConfig()
			redisCfg.SetParent(c.Alias).SetAutoLevel()
			outputCfg = redisCfg
		}
		ok, fn := manager.GetCreator(outputType)
		if ok {
			outputCores := fn(outputCfg)
			cores = append(cores, outputCores...)
		}
	}
	if len(cores) <= 0 {
		panic("logger no output,please register logger output")
	}
	zapLogger := zap.New(
		zap.NewTee(cores...),
		zapOptions...,
	)

	lg := &Logger{
		desugar:       zapLogger,
		config:        &c,
		sugar:         zapLogger.Sugar(),
		loggerType:    c.LoggerType,
		isPrintCommon: true,
	}
	return lg
}
