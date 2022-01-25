package database

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/adapter/mysql"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/database/config"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	userID  int64
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
	return "user"
}

func newMysqlAdapter() config.IDataBaseAdapter {
	ctx := context.Background()
	configDB := mysql.NewMysqlAdapter(ctx, &config.DBConfig{
		Type:            "mysql",
		DSN:             "dev:X#dSZ0PG*B@tcp(10.0.12.157:3306)/tabby-demo?charset=utf8mb4&interpolateParams=true",
		AutoConnect:     false,
		Debug:           true,
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

func Test_database_Count(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		sql    string
		values []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "test no failed",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "select id from #TABLE# ",
				values: nil,
			},
			want: 1,
		},

		{
			name: "test right",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "select count(*) from #TABLE# ",
				values: nil,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			if got := d.Count(tt.args.c, tt.args.sql, tt.args.values...); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_database_Create(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		dest   config.IModelConfig
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "测试创建指定字段",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx: context.Background(),
				dest: &User{
					Name:    "test",
					Address: "test",
				},
				sql:    "INSERT INTO `user` (`id`,`name`,`address`) VALUES (?,?,?)",
				values: []interface{}{2, "t2", "a2"},
			},
			wantErr: false,
		},
		{
			name: "测试创建指定字段2",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx: context.Background(),
				dest: &User{
					Name:    "test",
					Address: "test2",
				},
				sql:    "INSERT INTO #TABLE# (`id`,`name`,`address`) VALUES (?,?,?)",
				values: []interface{}{3, "t3", "a3"},
			},
			wantErr: false,
		},
		{
			name: "测试创建内容",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx: context.Background(),
				dest: &User{
					Name:    "test",
					Address: "test2",
				},
				sql:    "INSERT INTO #TABLE# (`id`,`name`,`address`) VALUES (?,?,?)",
				values: []interface{}{5, "t4", "a4"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			if err := d.Create(tt.args.dest, tt.args.dest, tt.args.sql, tt.args.values...); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_database_CreateOrUpdate(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test create or update",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "INSERT INTO `user` (`id`,`name`,`address`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `name`=?,`address`=?",
				values: []interface{}{2, "2", "3", "t1", "t2"},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "test create or update",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "INSERT INTO #TABLE# (`id`,`name`,`address`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `name`=?,`address`=?",
				values: []interface{}{2, "2", "3", "t1", "t23"},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			got, err := d.CreateOrUpdate(tt.args.c, tt.args.sql, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateOrUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_database_Delete(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test delete",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "",
				values: []interface{}{1},
			},
			want:    10,
			wantErr: true,
		},
		{
			name: "test delete",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "id =?",
				values: []interface{}{1},
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "test delete in ",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "id in ?",
				values: []interface{}{[]int64{1, 2, 131}},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "test delete = ",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "Delete  from #TABLE# ",
				values: nil,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			got, err := d.Delete(tt.args.c, tt.args.sql, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_database_GetField(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	var name string
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		dest   interface{}
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test filelds",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				dest:   &name,
				sql:    "select name from #TABLE# where id = ?",
				values: []interface{}{133},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			if err := d.GetField(tt.args.c, tt.args.dest, tt.args.sql, tt.args.values...); (err != nil) != tt.wantErr {
				t.Errorf("GetField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.Equal(t, "test", name)
}

func Test_database_GetOne(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	var res User
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		dest   interface{}
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test get one ",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				dest:   &res,
				sql:    "select * from #TABLE# where id =?",
				values: []interface{}{133},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			if err := d.GetOne(tt.args.c, tt.args.dest, tt.args.sql, tt.args.values...); (err != nil) != tt.wantErr {
				t.Errorf("GetOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.Equal(t, "test", res.Name)
	assert.Equal(t, "testall", res.Address)
}

func Test_database_Select(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	var ress []User
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx      context.Context
		c        config.IModelConfig
		destList interface{}
		sql      string
		values   []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test select",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:      context.Background(),
				c:        &User{},
				destList: &ress,
				sql:      "select * from #TABLE# where name =?",
				values:   []interface{}{"test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			if err := d.Select(tt.args.c, tt.args.destList, tt.args.sql, tt.args.values...); (err != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.NotEmpty(t, ress)
}

func Test_database_Update(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	type fields struct {
		forceMaster bool
		db          config.IDataBaseAdapter
	}
	type args struct {
		ctx    context.Context
		c      config.IModelConfig
		sql    string
		values []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test update",
			fields: fields{
				forceMaster: false,
				db:          mysqlDB,
			},
			args: args{
				ctx:    context.Background(),
				c:      &User{},
				sql:    "update #TABLE# set name=? where id =?",
				values: []interface{}{"t", 133},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				forceMaster: tt.fields.forceMaster,
				db:          tt.fields.db,
			}
			got, err := d.Update(tt.args.c, tt.args.sql, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_database_CreateBatch(t *testing.T) {
	mysqlDB := newMysqlAdapter()
	defer mysqlDB.Close()
	//ctx := context.Background()
	count := 10
	sql := "INSERT INTO `user` (`name`,`address`) VALUES "
	var values []interface{}
	for i := 0; i < int(count); i++ {
		sql += "(?, ?),"
		values = append(values, "name_"+strconv.Itoa(i), "add_"+strconv.Itoa(i))
	}
	sql = sql[0 : len(sql)-1]
	total, _ := mysqlDB.CreateBatch(sql, values...)
	assert.EqualValues(t, count, total)
}
