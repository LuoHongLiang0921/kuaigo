// @description gorm配置文件

package gormv2

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"


	"gorm.io/gorm/schema"

	"github.com/mitchellh/mapstructure"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Build 构建接口
type Build interface {
	Build(ctx context.Context) *DB
}

// StdConfig Deprecated
// StdConfig 标准配置 "dbs.{name}",其中name为别名
// 	@Description 标准配置 "dbs.{name}",其中name为别名
// 	@Param name key 名字
// 	@Return *DBConfig 载入后的配置
func StdConfig(ctx context.Context, name string) *DBConfig {
	return RawConfig(ctx, constant.DBS+name)
}

// RawConfig
// 	@Description 获取 db配置信息
//	@Param key 完整key 名，key 之间使用"."隔开
// 	@Return *DBConfig 载入后的配置
func RawConfig(ctx context.Context, key string) *DBConfig {
	config := DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
	}
	config.Name = key
	return config
}

// DataSources 自动识别多数据来源
type DataSources struct {
	// DBConfigs 所有配置， key为别名，value 为具体
	DBConfigs map[string]*DBConfig
}

// GetDefaultConfig
// @Description: 获取多数据源默认配置
// @receiver ds DataSources
// @Return *DBConfig 载入后的多数据源配置
func (ds *DataSources) GetDefaultConfig() *DBConfig {
	if ds.DBConfigs == nil {
		return nil
	}
	return ds.DBConfigs[constant.AlisDefault]
}

// DataSourcesConfig 获取多数据源配置信息
// @Description 获取多数据源配置信息，根据表名自动识别数据源
// @Param ctx 上下文
// @Param key dbs 根 key 名字
// @Return *DataSources 多个数据源配置
func DataSourcesConfig(ctx context.Context, key string) *DataSources {
	dataV := conf.GetStringMap(key)
	m := make(map[string]*DBConfig)
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &m,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	err = decoder.Decode(dataV)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	return &DataSources{
		DBConfigs: m,
	}
}

// MultiDB
// 	@Description: 载入配置并实例化多数据源db
//	@Param ctx 上下文
//	@Param key 根key 名字
// 	@Return map[string]*gorm.DB 实例化后的多数据源键值对，key 为命名，value 为 gorm DB 实例
func MultiDB(ctx context.Context, key string) map[string]*gorm.DB {
	dataV := conf.GetStringMap(key)
	m := make(map[string]*DBConfig)
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &m,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	err = decoder.Decode(dataV)
	if err != nil {
		klog.KuaigoLogger.Panic("unmarshal key", klog.FieldMod("gorm"), klog.FieldErr(err), klog.FieldKey(key))
		return nil
	}
	dbMap := make(map[string]*gorm.DB)
	for k, v := range m {
		db := v.Build(ctx)
		dbMap[k] = db
	}
	return dbMap
}

// Build 多数据源构建
// @Description  DB 解析器,多个数据库支持。数据源切换规则，主从分离
// @receiver cfg
// @Param ctx 上下文信息
// @Return *DB 数据库实例
func (ds *DataSources) Build(ctx context.Context) *DB {
	// 取 default
	defaultCfg := ds.GetDefaultConfig()
	if defaultCfg == nil {
		klog.KuaigoLogger.Panic("多数据需要默认配置(mysql.default)", klog.FieldMod("gormv2"))
	}
	rootDB, err := gorm.Open(mysql.Open(defaultCfg.DSN), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})
	if err != nil {
		klog.KuaigoLogger.Panicf("init root database err:%v", err)
	}
	rs := new(dbresolver.DBResolver)
	for _, rawCfg := range ds.DBConfigs {
		var dates []interface{}
		for _, v := range rawCfg.Datas {
			dates = append(dates, v)
		}
		dbCfg := rawCfg.toDBResolver()
		rs = rs.Register(dbCfg, dates...).SetConnMaxIdleTime(rawCfg.DialTimeout).
			SetConnMaxLifetime(rawCfg.ConnMaxLifetime).
			SetMaxIdleConns(rawCfg.MaxIdleConns).
			SetMaxOpenConns(rawCfg.MaxOpenConns)
	}

	err = rootDB.Use(rs)
	if err != nil {
		klog.KuaigoLogger.Panicf("init database err:%v", err)
	}
	return rootDB
}

// DBConfig 数据库配置
type DBConfig struct {
	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" yaml:"dsn" `
	//Sources 数据源
	Sources []string `json:"sources" yaml:"sources"`
	//Replicas 从数据源
	Replicas []string `json:"replicas" yaml:"replicas"`
	// Datas 数据
	Datas []string `json:"datas" yaml:"datas"`
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
}

// DefaultConfig
// 	@Description 获取单数据源默认配置
// 	@Return *DBConfig 设置默认值后的DB配置
func DefaultConfig() *DBConfig {
	return &DBConfig{
		DSN:             "",
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
	}
}

// WithLogger
// 	@Description 设置日志实例
// 	@Receiver config DBConfig
//	@Param log 日志实例
// 	@Return *Config 设置日志后的配置
func (cfg *DBConfig) WithLogger(log *klog.Logger) *DBConfig {
	cfg.logger = log
	return cfg
}

func (cfg *DBConfig) toDBResolver() dbresolver.Config {
	var sources, replicas []gorm.Dialector
	for _, v := range cfg.Sources {
		sources = append(sources, mysql.Open(v))
	}
	for _, v := range cfg.Replicas {
		replicas = append(replicas, mysql.Open(v))
	}
	var policy dbresolver.Policy
	if len(replicas) > 0 {
		policy = dbresolver.RandomPolicy{}
	}

	return dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		Policy:   policy,
	}
}

// Build
// @Description 单源配置文件
// @receiver config DBConfig
// @Param ctx 上下文
// @Return *DB 初始化后的单源DB
func (cfg *DBConfig) Build(ctx context.Context) *DB {
	var gormLogger logger.Interface
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		DryRun: cfg.DryRun,
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})
	if err != nil {
		if cfg.OnDialError == "panic" {
			klog.KuaigoLogger.Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		} else {
			klog.KuaigoLogger.Info(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		}
	}

	inner, err := db.DB()
	if err != nil {
		if cfg.OnDialError == "panic" {
			klog.KuaigoLogger.Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		} else {
			klog.KuaigoLogger.Info(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		}
	}
	// 设置默认连接配置
	inner.SetMaxIdleConns(cfg.MaxIdleConns)
	inner.SetMaxOpenConns(cfg.MaxOpenConns)

	if cfg.ConnMaxLifetime != 0 {
		inner.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	return db
}
