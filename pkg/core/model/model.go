package model

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
)

type BaseModel struct {
	config string
	ctx    context.Context
	db     config.IDatabase
}

// Build
//  @Description 构造函数，New完之后必须调用。初始化前置资源
func (m *BaseModel) Build(ctx context.Context) *BaseModel {
	m.ctx = ctx
	m.db = database.GetDatabaseFactoryInstance().GetDatabase(ctx, m.config)
	return m
}

// GetDb
//  @Description 获取数据库操作
// 	@Receiver BaseModel
// 	@Return IDatabase 数据库操作接口
func (m *BaseModel) GetDb() config.IDatabase {
	return m.db
}

// WithConfig
// 	@Description 数据库配置
// 	@Receiver BaseModel
//	@Param config 数据库配置key
// 	@Return *BaseModel
func (m *BaseModel) WithConfig(config string) *BaseModel {
	m.config = config
	return m
}
