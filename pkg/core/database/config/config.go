// @Description 框架全局公共属性

package config

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"

)

const databaseConfigPrefix = "dbs."

// DBConfig 数据库配置
type DBConfig struct {
	Name string
	Type string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" yaml:"dsn" `
	//Sources 数据源
	Sources []string `json:"sources" yaml:"sources"`
	//Replicas 从数据源
	Replicas []string `json:"replicas" yaml:"replicas"`
	// Datas 数据
	Datas []string `json:"datas" yaml:"datas"`
	//启动是否自动连接
	AutoConnect bool `json:"autoConnect" yaml:"autoConnect"`
	// Debug开关
	Debug bool `json:"debug" yaml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" yaml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" yaml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" yaml:"connMaxLifetime"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
	// 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	// 关闭指标采集
	DisableMetric bool `json:"disableMetric" yaml:"disableMetric"`
	// 关闭链路追踪
	DisableTrace bool `json:"disableTrace" yaml:"disableTrace"`
	// 记录错误sql时,是否打印包含参数的完整sql语句
	// select * from aid = ?;
	// select * from aid = 288016;
	DetailSQL bool `json:"detailSql" yaml:"detailSql"`
	DryRun    bool
	raw       interface{}
	logger    *klog.Logger

	change    chan struct{}
	latestDsn string
}

// GetConfig
// 	@Description 获取db配置信息
//	@Param key 完整key 名，key 之间使用"."隔开
// 	@Return *DBConfig 载入后的配置
func GetConfig(key string) *DBConfig {
	config := &DBConfig{
		Type:            "mysql",
		DSN:             "",
		AutoConnect:     false,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: ktime.Duration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   ktime.Duration("500ms"),
		DialTimeout:     ktime.Duration("1s"),
		DisableMetric:   false,
		DisableTrace:    false,
		raw:             nil,
		logger:          klog.KuaigoLogger,
		change:          make(chan struct{}, 1),
	}
	config.latestDsn = config.DSN
	configKey := config.getConfigKey(key)
	if err := conf.UnmarshalKey(configKey, config); err != nil {
		klog.KuaigoLogger.Panicf("unmarshal key %v err %v", configKey, err)
	}
	config.latestDsn = config.DSN
	config.Name = key
	config.setOnChange(key)
	return config
}

// IsConfigChange
// 	@Description
// 	@Receiver c
// 	@Return <-chan
func (c *DBConfig) IsConfigChange() <-chan struct{} {
	return c.change
}

func (c *DBConfig) getConfigKey(key string) string {
	return databaseConfigPrefix + key
}

func (c *DBConfig) setOnChange(key string) {
	configKey := c.getConfigKey(key)
	conf.OnChange(func(cfg *conf.Configuration) {
		dsnRoot := configKey + ".dsn"
		dsnStr := cfg.GetString(dsnRoot)
		klog.Debugf("%s change, result %ss", dsnRoot, dsnStr)
		if dsnStr != "" && dsnStr != c.latestDsn {
			c.change <- struct{}{}
		}
		c.latestDsn = dsnStr
	})
}
