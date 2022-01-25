// @description gorm v2 例子

package gormv2_test

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/storage/gormv2"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

type User struct {
	ID   int64
	Name string
}

func (u User) TableName() string {
	return "user"
}

type Device struct {
	ID   int64
	UDID string
}

func (d Device) TableName() string {
	return "device"
}

// 单数据源
func ExampleDBConfig_Build() {
	ctx := context.Background()
	db := gormv2.RawConfig(ctx, "mysql").Build(ctx)
	var user User
	err := db.Select("id", "name").Where("id =?", 2).First(&user).Error
	if err != nil {
		klog.WithContext(ctx).Errorf("get user id %v", 2)
		return
	}
	klog.WithContext(ctx).Infof("user %+v", user)
}

// 多数据源
func ExampleDataSources_Build() {
	ctx := context.Background()
	multiDB := gormv2.DataSourcesConfig(ctx, "mysql").Build(ctx)
	// 查询user
	var user User
	err := multiDB.Select("id", "name").Where("id =?", 2).First(&user).Error
	if err != nil {
		klog.WithContext(ctx).Errorf("get user id %v", 2)
		return
	}
	klog.WithContext(ctx).Infof("user %+v", user)
	// 查询 device
	var device Device
	err = multiDB.Select("id", "udid").Where("id =?", 2).First(&device).Error
	if err != nil {
		klog.WithContext(ctx).Errorf("get user id %v", 2)
		return
	}
	klog.WithContext(ctx).Infof("user %+v", user)
}
