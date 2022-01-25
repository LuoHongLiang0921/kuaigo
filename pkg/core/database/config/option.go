package config

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type (
	DB      = gorm.DB
	Config  = gorm.Config
	Session = gorm.Session
)

var (
	errSlowCommand    = errors.New("mysql slow command")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// IModelConfig 数据模型需要实现接口
type IModelConfig interface {
	TableName() string

	PrimaryKey() string

	CachePrefix() string
}

// IDatabase 数据操作实现接口
type IDatabase interface {
	GetRawDB() *DB
	WithContext(ctx context.Context) IDatabase
	ForceMaster(forceMaster bool)
	Count(c IModelConfig, sql string, values ...interface{}) int64
	GetField(c IModelConfig, dest interface{}, sql string, values ...interface{}) error
	GetOne(c IModelConfig, dest interface{}, sql string, values ...interface{}) error
	Select(c IModelConfig, destList interface{}, sql string, values ...interface{}) error
	Create(c IModelConfig, dest interface{}, sql string, values ...interface{}) error
	CreateBatch(c IModelConfig, sql string, values ...interface{}) (int64, error)
	Update(c IModelConfig, sql string, values ...interface{}) (int64, error)
	CreateOrUpdate(c IModelConfig, sql string, values ...interface{}) (int64, error)
	Delete(c IModelConfig, sql string, values ...interface{}) (int64, error)
}

// IDataBaseAdapter 数据库适配器接口
type IDataBaseAdapter interface {
	WithContext(ctx context.Context) IDataBaseAdapter
	Count(sql string, values ...interface{}) int64
	GetField(dest interface{}, sql string, values ...interface{}) error
	GetOne(dest interface{}, sql string, values ...interface{}) error
	Select(destList interface{}, sql string, values ...interface{}) error
	Create(dest interface{}, sql string, values ...interface{}) error
	CreateBatch(sql string, values ...interface{}) (int64, error)
	Update(sql string, values ...interface{}) (int64, error)
	CreateOrUpdate(sql string, values ...interface{}) (int64, error)
	Delete(sql string, values ...interface{}) (int64, error)
	Open(ctx context.Context) *DB
	GetDB() *DB
	Close()
}
