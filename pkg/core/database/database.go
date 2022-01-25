// @Description 数据库操作类

package database

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
	"strings"

)

type database struct {
	forceMaster bool
	db          config.IDataBaseAdapter
	ctx         context.Context
}

// NewDatabase
// 	@Description 数据库操作类构建函数
//  @Param ctx 上下文Context
//	@Param conf 数据库配置key
// 	@Return IDatabase
func NewDatabase(dataBaseAdapter config.IDataBaseAdapter) config.IDatabase {
	return &database{
		forceMaster: false,
		db:          dataBaseAdapter,
	}
}

// WithContext
// 	@Description
// 	@Receiver d
//	@Param ctx
// 	@Return config.IDatabase
func (d *database) WithContext(ctx context.Context) config.IDatabase {
	if ctx == nil {
		return d
	}
	newR := d.clone()
	newR.ctx = ctx
	newR.db = d.db.WithContext(ctx)
	return newR
}

func (d *database) clone() *database {
	copy := *d
	return &copy
}

// GetRawDB
// 	@Description 获取最底层数据库操作句柄,慎用
// 	@Receiver Database
// 	@Return config.DB 底层数据库操作类
func (d *database) GetRawDB() *config.DB {
	return d.db.GetDB()
}

// ForceMaster
// 	@Description 强制使用主库
// 	@Receiver Database
//  @Param forceMaster 是否强制使用主库
func (d *database) ForceMaster(forceMaster bool) {
	d.forceMaster = forceMaster
}

// Count
// 	@Description 计算条数,是GetField函数的包装
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return int64 条数
func (d *database) Count(c config.IModelConfig, sql string, values ...interface{}) int64 {
	var total int64 = 0
	err := d.db.GetField(&total, d.parseTableName(c, sql), values...)
	if err != nil {
		return 0
	}
	return total
}

// GetField
// 	@Description 获取单个字段值
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param dest 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (d *database) GetField(c config.IModelConfig, dest interface{}, sql string, values ...interface{}) error {
	return d.db.GetField(dest, d.parseTableName(c, sql), values...)
}

// GetOne
// 	@Description 获取单条数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param dest 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (d *database) GetOne(c config.IModelConfig, dest interface{}, sql string, values ...interface{}) error {
	return d.db.GetOne(dest, d.parseTableName(c, sql), values...)
}

// Select
// 	@Description 获取单条数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param destList 传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (d *database) Select(c config.IModelConfig, destList interface{}, sql string, values ...interface{}) error {
	return d.db.Select(destList, d.parseTableName(c, sql), values...)
}

// Create
// 	@Description 创建数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param dest 创建成功后传入的接收结果数据地址
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 错误
func (d *database) Create(c config.IModelConfig, dest interface{}, sql string, values ...interface{}) error {
	return d.db.Create(dest, d.parseTableName(c, sql), values...)
}

// CreateBatch
// 	@Description 创建数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return 写入条数，错误
func (d *database) CreateBatch(c config.IModelConfig, sql string, values ...interface{}) (int64, error) {
	return d.db.CreateBatch(d.parseTableName(c, sql), values...)
}

// Update
// 	@Description 更新数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (d *database) Update(c config.IModelConfig, sql string, values ...interface{}) (int64, error) {
	return d.db.Update(d.parseTableName(c, sql), values...)
}

// CreateOrUpdate
// 	@Description 创建或更新数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param sql SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (d *database) CreateOrUpdate(c config.IModelConfig, sql string, values ...interface{}) (int64, error) {
	return d.db.CreateOrUpdate(d.parseTableName(c, sql), values...)
}

// Delete
// 	@Description 删除数据
// 	@Receiver Database
//  @Param c 数据表配置接口
//  @Param where where SQL语句模板
//  @Param values SQL语句参数
// 	@Return error 影响条数和错误
func (d *database) Delete(c config.IModelConfig, sql string, values ...interface{}) (int64, error) {
	return d.db.Delete(d.parseTableName(c, sql), values...)
}

// parseTableName
// 	@Description 内部解析数据表名函数
//  @Param c 数据表配置接口
//  @Param sql SQL语句模板
// 	@Return string 更改表名后的sql模板
func (d *database) parseTableName(c config.IModelConfig, sql string) string {
	return strings.ReplaceAll(sql, "#TABLE#", "`"+c.TableName()+"`")
}
