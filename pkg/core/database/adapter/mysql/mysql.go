package mysql

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type mysqlAdapter struct {
	ctx       context.Context
	config    *config.DBConfig
	db        *config.DB
	openOnce  sync.Once
	hasOpen   bool
	openError error

	mux sync.RWMutex
}

// NewMysqlAdapter
// 	@Description Mysql数据库操作类构建函数
//  @Param ctx 上下文Context
//	@Param config 数据库配置
// 	@Return IDataBaseAdapter
func NewMysqlAdapter(ctx context.Context, config *config.DBConfig) config.IDataBaseAdapter {
	ret := &mysqlAdapter{ctx: ctx,
		config:    config,
		hasOpen:   false,
		openError: errors.New("database " + config.Name + " open error"),
	}
	if config.AutoConnect {
		ret.Open(ctx)
	}
	ret.onDsnChange()
	return ret
}

func (m *mysqlAdapter) onDsnChange() {
	kgo.SafeGo(func() {
		for range m.config.IsConfigChange() {
			// todo: renew db 句柄
			// o句柄ld 句柄怎么处理，不接受新的请求，只接受新的
			m.setDB(m.ctx)
		}

	}, func(err error) {
		klog.Warnf("gorm change err:%v", err)
	})
}

func (m *mysqlAdapter) WithContext(ctx context.Context) config.IDataBaseAdapter {
	if ctx == nil {
		return m
	}
	newR := m.clone()
	newR.ctx = ctx
	return newR
}

func (m *mysqlAdapter) clone() *mysqlAdapter {
	copy := mysqlAdapter{
		ctx:       m.ctx,
		config:    m.config,
		db:        m.db,
		openOnce:  sync.Once{},
		hasOpen:   m.hasOpen,
		openError: m.openError,
		mux:       sync.RWMutex{},
	}
	return &copy
}

// Open
// 	@Description 打开数据库连接
// 	@Receiver mysqlAdapter
// 	@Return config.DB 底层数据库操作类
func (m *mysqlAdapter) Open(ctx context.Context) *config.DB {
	m.openOnce.Do(func() {
		m.setDB(ctx)
	})
	return m.db
}

func (m *mysqlAdapter) setDB(ctx context.Context) {
	var gormLogger logger.Interface
	if m.config.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(m.config.DSN), &gorm.Config{
		DryRun: m.config.DryRun,
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})
	if err != nil {
		if m.config.OnDialError == "panic" {
			klog.KuaigoLogger.WithContext(ctx).Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		} else {
			klog.KuaigoLogger.WithContext(ctx).Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		}
	}

	rs := new(dbresolver.DBResolver)
	var dates []interface{}
	for _, v := range m.config.Datas {
		dates = append(dates, v)
	}
	dbCfg := m.parserDBResolver()
	rs = rs.Register(dbCfg, dates...).SetConnMaxIdleTime(m.config.DialTimeout).
		SetConnMaxLifetime(m.config.ConnMaxLifetime).
		SetMaxIdleConns(m.config.MaxIdleConns).
		SetMaxOpenConns(m.config.MaxOpenConns)

	err = db.Use(rs)
	if err != nil {
		if m.config.OnDialError == "panic" {
			klog.KuaigoLogger.WithContext(ctx).Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		} else {
			klog.KuaigoLogger.WithContext(ctx).Panicf(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		}
	}
	inner, err := db.DB()
	if err != nil {
		if m.config.OnDialError == "panic" {
			klog.KuaigoLogger.WithContext(ctx).Panic(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		} else {
			klog.KuaigoLogger.WithContext(ctx).Panicf(ecode.MsgClientMysqlOpenStart, klog.FieldMod("gorm"), klog.FieldErr(err))
		}
	}
	// 设置默认连接配置
	inner.SetConnMaxIdleTime(m.config.DialTimeout)
	inner.SetConnMaxLifetime(m.config.ConnMaxLifetime)
	inner.SetMaxIdleConns(m.config.MaxIdleConns)
	inner.SetMaxOpenConns(m.config.MaxOpenConns)

	m.mux.Lock()
	m.db = db
	m.hasOpen = true
	m.mux.Unlock()
}

// GetDB
// 	@Description 获取已打开数据库
// 	@Receiver mysqlAdapter
//  @Param ctx 上下文Context
// 	@Return config.DB 底层数据库操作类
func (m *mysqlAdapter) GetDB() *config.DB {
	if m.checkOpen() {
		return m.getDB()
	}
	return m.getDB()
}

func (m *mysqlAdapter) getDB() *config.DB {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.db
}

// Close
// 	@Description 关闭已打开数据库
// 	@Receiver mysqlAdapter
func (m *mysqlAdapter) Close() {
	if m.hasOpen {
		db, err := m.getDB().DB()
		if err == nil && db != nil {
			db.Close()
		}
	}
}

// Count
// 	@Description 计算条数,是GetField函数的包装
// 	@Receiver mysqlAdapter
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return int64 条数
func (m *mysqlAdapter) Count(sql string, values ...interface{}) int64 {
	var total int64 = 0
	err := m.GetField(&total, sql, values...)
	if err != nil {
		return 0
	}
	return total
}

// GetField
// 	@Description 获取单个字段值
// 	@Receiver mysqlAdapter
//  @Param dest 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (m *mysqlAdapter) GetField(dest interface{}, sql string, values ...interface{}) error {
	if m.checkOpen() {
		return m.getDB().WithContext(m.getContext()).Raw(sql, values...).Scan(dest).Error
	}
	return m.openError
}

func (m *mysqlAdapter) getContext() context.Context {
	if m.ctx == nil {
		return context.Background()
	}
	return m.ctx
}

// GetOne
// 	@Description 获取单条数据
// 	@Receiver mysqlAdapter
//  @Param dest 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (m *mysqlAdapter) GetOne(dest interface{}, sql string, values ...interface{}) error {
	if m.checkOpen() {
		return m.getDB().WithContext(m.getContext()).Raw(sql, values...).Take(dest).Error
	}
	return m.openError
}

// Select
// 	@Description 获取单条数据
// 	@Receiver Database
//  @Param destList 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (m *mysqlAdapter) Select(destList interface{}, sql string, values ...interface{}) error {
	if m.checkOpen() {
		return m.getDB().WithContext(m.getContext()).Raw(sql, values...).Find(destList).Error
	}
	return m.openError
}

// Create
// 	@Description 创建数据
// 	@Receiver mysqlAdapter
//  @Param dest 创建成功后传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (m *mysqlAdapter) Create(dest interface{}, sql string, values ...interface{}) error {
	if m.checkOpen() {
		return m.getDB().WithContext(m.getContext()).Raw(sql, values...).Create(dest).Error
	}
	return m.openError
}

// CreateBatch
// 	@Description 创建数据
// 	@Receiver mysqlAdapter
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return 写入条数，错误
func (m *mysqlAdapter) CreateBatch(sql string, values ...interface{}) (int64, error) {
	if m.checkOpen() {
		db := m.getDB().WithContext(m.getContext()).Exec(sql, values...)
		return db.RowsAffected, db.Error
	}
	return 0, m.openError
}

// Update
// 	@Description 更新数据
// 	@Receiver mysqlAdapter
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (m *mysqlAdapter) Update(sql string, values ...interface{}) (int64, error) {
	if m.checkOpen() {
		db := m.getDB().WithContext(m.getContext()).Exec(sql, values...)
		return db.RowsAffected, db.Error
	}
	return 0, m.openError
}

// CreateOrUpdate
// 	@Description 创建或更新数据
// 	@Receiver mysqlAdapter
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (m *mysqlAdapter) CreateOrUpdate(sql string, values ...interface{}) (int64, error) {
	if m.checkOpen() {
		db := m.getDB().WithContext(m.getContext()).Exec(sql, values...)
		return db.RowsAffected, db.Error
	}
	return 0, m.openError
}

// Delete
// 	@Description 删除数据
// 	@Receiver mysqlAdapter
//  @Param where where条件语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (m *mysqlAdapter) Delete(sql string, values ...interface{}) (int64, error) {
	if m.checkOpen() {
		db := m.getDB().WithContext(m.getContext()).Exec(sql, values...)
		return db.RowsAffected, db.Error
	}
	return 0, m.openError
}

// checkOpen
// 	@Description 检查连接是否打开
// 	@Receiver mysqlAdapter
// 	@Return 打开结果
func (m *mysqlAdapter) checkOpen() bool {
	if !m.hasOpen {
		m.Open(m.ctx)
	}
	return true
}

// parserDBResolver
// 	@Description 增加数据主从支持
// 	@Receiver mysqlAdapter
// 	@Return resolver.Config 构造好的主从配置
func (m *mysqlAdapter) parserDBResolver() dbresolver.Config {
	var sources, replicas []gorm.Dialector
	for _, v := range m.config.Sources {
		sources = append(sources, mysql.Open(v))
	}
	for _, v := range m.config.Replicas {
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
