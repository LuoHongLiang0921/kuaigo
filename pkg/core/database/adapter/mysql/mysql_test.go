package mysql

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type User struct {
	ID      int64
	Name    string
	Address string
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) PrimaryKey() string {
	return "id"
}

func (u *User) CachePrefix() string {
	return "id"
}

func newMysqlAdapter() config.IDataBaseAdapter {
	ctx := context.Background()
	configDB := NewMysqlAdapter(ctx, &config.DBConfig{
		Type:            "mysql",
		DSN:             "dev:X#dSZ0PG*B@tcp(10.0.12.157:3306)/tabby-demo?charset=utf8mb4&interpolateParams=true",
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
		DryRun:          false,
	})
	return configDB
}

func Test_mysqlAdaper_Create(t *testing.T) {
	cfgDB := newMysqlAdapter()
	defer cfgDB.Close()
	user := &User{}
	err := cfgDB.Create(&user, "INSERT INTO `user` (`name`,`address`) VALUES (?,?)", "name", "address")

	assert.NoError(t, err, "test error")
}

func Test_mysqlAdapter_GetField(t *testing.T) {
	cfgDB := newMysqlAdapter()
	defer cfgDB.Close()
	var total int64
	err := cfgDB.GetField(&total, "select count(*) from user")
	assert.NoError(t, err, "not have user")
	assert.Equal(t, int64(1), total)
}

func Test_mysqlAdapter_Count(t *testing.T) {
	cfgDB := newMysqlAdapter()
	defer cfgDB.Close()
	//var total int64
	total := cfgDB.Count("select count(*) from user")
	//assert.NoError(t, err, "not have user")
	assert.Equal(t, int64(1), total)
}

func Test_mysqlAdapter_CreateBatch(t *testing.T) {
	cfgDB := newMysqlAdapter()
	defer cfgDB.Close()
	count := 10
	sql := "INSERT INTO `user` (`name`,`address`) VALUES "
	var values []interface{}
	for i := 0; i < int(count); i++ {
		sql += "(?, ?),"
		values = append(values, "name_"+strconv.Itoa(i), "add_"+strconv.Itoa(i))
	}
	sql = sql[0 : len(sql)-1]
	total, _ := cfgDB.CreateBatch(sql, values...)
	assert.Equal(t, count, total)
}
