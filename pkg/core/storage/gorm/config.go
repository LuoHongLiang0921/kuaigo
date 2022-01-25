// @Description

package gorm

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/metric"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"

)

// StdConfig 标准配置，规范配置文件头
func StdConfig(name string) *Config {
	return RawConfig("tabby.mysql." + name)
}

// RawConfig 传入mapstructure格式的配置
// example: RawConfig("tabby.mysql.stt_config")
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config, conf.TagName("toml")); err != nil {
		klog.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
	}
	config.Name = key
	return config
}

// config options
type Config struct {
	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug开关
	Debug bool `json:"debug" toml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	// 关闭指标采集
	DisableMetric bool `json:"disableMetric" toml:"disableMetric"`
	// 关闭链路追踪
	DisableTrace bool `json:"disableTrace" toml:"disableTrace"`

	// 记录错误sql时,是否打印包含参数的完整sql语句
	// select * from aid = ?;
	// select * from aid = 288016;
	DetailSQL bool `json:"detailSql" toml:"detailSql"`

	raw          interface{}
	logger       *klog.Logger
	interceptors []Interceptor
	dsnCfg       *DSN
}

// DefaultConfig
// 	@Description 返回默认配置
// 	@Return *Config
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: ktime.Duration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   ktime.Duration("500ms"),
		DialTimeout:     ktime.Duration("1s"),
		DisableMetric:   true,
		DisableTrace:    true,
		raw:             nil,
		logger:          klog.KuaigoLogger,
	}
}

// WithLogger
// 	@Description 设置日志实例
// 	@Receiver config Config
//	@Param log 日志实例
// 	@Return *Config 设置日志后的配置
func (config *Config) WithLogger(log *klog.Logger) *Config {
	config.logger = log
	return config
}

// WithInterceptor
// 	@Description 设置拦截器
// 	@Receiver config Config
//	@Param intes 拦截器数组
// 	@Return *Config 设置拦截器数组后的配置
func (config *Config) WithInterceptor(intes ...Interceptor) *Config {
	if config.interceptors == nil {
		config.interceptors = make([]Interceptor, 0)
	}
	config.interceptors = append(config.interceptors, intes...)
	return config
}

func (config *Config) getContext() context.Context {
	return context.TODO()
}

// Build
// 	@Description 构造 gorm db 实例
// 	@Receiver config 配置
// 	@Return *DB 实例后的gorm db 实例
func (config *Config) Build() *DB {
	var err error
	config.dsnCfg, err = ParseDSN(config.DSN)
	if err == nil {
		config.logger.Info(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldAddr(config.dsnCfg.Addr), klog.FieldName(config.dsnCfg.DBName))
	} else {
		config.logger.Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
	}

	if config.Debug {
		config = config.WithInterceptor(debugInterceptor)
	}

	if !config.DisableMetric {
		config = config.WithInterceptor(metricInterceptor)
	}

	db, err := Open(config.getContext(), "mysql", config)
	if err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("open mysql", klog.FieldMod("gorm"), klog.FieldErrKind(ecode.ErrKindRequestErr), klog.FieldErr(err), klog.FieldAddr(config.dsnCfg.Addr), klog.FieldValueAny(config))
		} else {
			metric.LibHandleCounter.Inc(metric.TypeGorm, config.Name+".ping", config.dsnCfg.Addr, "open err")
			config.logger.Error("open mysql", klog.FieldMod("gorm"), klog.FieldErrKind(ecode.ErrKindRequestErr), klog.FieldErr(err), klog.FieldAddr(config.dsnCfg.Addr), klog.FieldValueAny(config))
			return db
		}
	}

	if err := db.DB().Ping(); err != nil {
		config.logger.Panic("ping mysql", klog.FieldMod("gorm"), klog.FieldErrKind(ecode.ErrKindRequestErr), klog.FieldErr(err), klog.FieldValueAny(config))
	}

	// storage db
	instances.Store(config.Name, db)
	return db
}
