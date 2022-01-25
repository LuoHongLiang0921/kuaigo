// @Description

package gorm

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kdebug"

	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// SQLCommon ...
type (
	// SQLCommon alias of gorm.SQLCommon
	SQLCommon = gorm.SQLCommon
	// Callback alias of gorm.Callback
	Callback = gorm.Callback
	// CallbackProcessor alias of gorm.CallbackProcessor
	CallbackProcessor = gorm.CallbackProcessor
	// Dialect alias of gorm.Dialect
	Dialect = gorm.Dialect
	// Scope ...
	Scope = gorm.Scope
	// DB ...
	DB = gorm.DB
	// Model ...
	Model = gorm.Model
	// ModelStruct ...
	ModelStruct = gorm.ModelStruct
	// Field ...
	Field = gorm.Field
	// FieldStruct ...
	StructField = gorm.StructField
	// RowQueryResult ...
	RowQueryResult = gorm.RowQueryResult
	// RowsQueryResult ...
	RowsQueryResult = gorm.RowsQueryResult
	// Association ...
	Association = gorm.Association
	// Errors ...
	Errors = gorm.Errors
	// logger ...
	Logger = gorm.Logger
)

var (
	errSlowCommand = errors.New("mysql slow command")

	// IsRecordNotFoundError ...
	IsRecordNotFoundError = gorm.IsRecordNotFoundError

	// ErrRecordNotFound returns a "record not found error". Occurs only when attempting to query the database with a struct; querying with a slice won't return this error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidSQL occurs when you attempt a query with invalid SQL
	ErrInvalidSQL = gorm.ErrInvalidSQL
	// ErrInvalidTransaction occurs when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = gorm.ErrUnaddressable
)

// WithContext
// 	@Description 设置上下文
//	@Param ctx 上下文
//	@Param db gorm db 实例
// 	@Return *DB 设置上下文后db 实例
func WithContext(ctx context.Context, db *DB) *DB {
	db.InstantSet("_context", ctx)
	return db
}

// Open
// 	@Description 初始化一个新的db 连接
//	@Param ctx 上下文
//	@Param dialect
//	@Param options 配置选项
// 	@Return *DB 新的db 连接
// 	@Return error 错误
func Open(ctx context.Context, dialect string, options *Config) (*DB, error) {
	inner, err := gorm.Open(dialect, options.DSN)
	if err != nil {
		return nil, err
	}

	inner.LogMode(options.Debug)
	// 设置默认连接配置
	inner.DB().SetMaxIdleConns(options.MaxIdleConns)
	inner.DB().SetMaxOpenConns(options.MaxOpenConns)

	if options.ConnMaxLifetime != 0 {
		inner.DB().SetConnMaxLifetime(options.ConnMaxLifetime)
	}

	if kdebug.IsDevelopmentMode() {
		inner.LogMode(true)
	}

	replace := func(ctx context.Context, processor func() *gorm.CallbackProcessor, callbackName string, interceptors ...Interceptor) {
		old := processor().Get(callbackName)
		var handler = old
		for _, inte := range interceptors {
			handler = inte(ctx, options.dsnCfg, callbackName, options)(handler)
		}
		processor().Replace(callbackName, handler)
	}

	replace(
		ctx,
		inner.Callback().Delete,
		"gorm:delete",
		options.interceptors...,
	)
	replace(
		ctx,
		inner.Callback().Update,
		"gorm:update",
		options.interceptors...,
	)
	replace(
		ctx,
		inner.Callback().Create,
		"gorm:create",
		options.interceptors...,
	)
	replace(
		ctx,
		inner.Callback().Query,
		"gorm:query",
		options.interceptors...,
	)
	replace(
		ctx,
		inner.Callback().RowQuery,
		"gorm:row_query",
		options.interceptors...,
	)

	return inner, err
}
