package database

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/adapter/mysql"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
	"sync"

)

var databaseFactoryInstance *databaseFactory
var instanceOnce sync.Once
var sum sync.RWMutex

type databaseFactory struct {
	dataBases map[string]config.IDatabase
}

// GetDatabaseFactoryInstance
//  @Description 获取全局单例，如果new多个，永远指向最后一次Build的实例
func GetDatabaseFactoryInstance() *databaseFactory {
	instanceOnce.Do(func() {
		databaseFactoryInstance = new(databaseFactory)
		databaseFactoryInstance.dataBases = make(map[string]config.IDatabase, 0)
	})
	return databaseFactoryInstance
}

// GetDatabase
// 	@Description 获取数据库操作类
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param conf 数据库配置key
// 	@Return IDataBaseAdapter
func (df *databaseFactory) GetDatabase(ctx context.Context, conf string) config.IDatabase {
	if db, ok := df.dataBases[conf]; ok {
		return db
	}
	return df.buildDatabase(ctx, conf)
}

// buildDatabase
// 	@Description 构造数据库内部方法
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param config 数据库配置key
// 	@Return IDatabase
func (df *databaseFactory) buildDatabase(ctx context.Context, conf string) config.IDatabase {
	sum.Lock()
	if _, ok := df.dataBases[conf]; !ok {
		df.dataBases[conf] = NewDatabase(df.buildDatabaseAdapter(ctx, conf))
	}
	sum.Unlock()
	return df.dataBases[conf]
}

// buildDatabaseDriver
// 	@Description 构造数据库适配器内部方法
// 	@Receiver databaseFactory
//  @Param ctx 上下文Context
//	@Param config 数据库配置key
// 	@Return IDataBaseAdapter
func (df *databaseFactory) buildDatabaseAdapter(ctx context.Context, conf string) config.IDataBaseAdapter {
	dbConfig := config.GetConfig(conf)
	if dbConfig.Type == "mysql" {
		return mysql.NewMysqlAdapter(ctx, dbConfig)
	}
	panic("tabby not support " + dbConfig.Type + " database")
}
