package database

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/file"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func initTest() {
	file.RegisterConfigHandler()
	configAddr := os.Getenv("CONFIG_FILE_ADDR")
	provider, err := manager.NewConfigSource(configAddr)
	if err == manager.ErrConfigAddr {
		klog.Panic(err.Error())
	}
	if err := conf.LoadFromConfigSource(provider, yaml.Unmarshal); err != nil {
		klog.Panic(err.Error())
	}
}

func TestGetDatabaseFactoryInstance(t *testing.T) {
	initTest()
	var wg sync.WaitGroup
	rootCtx := context.Background()
	db := GetDatabaseFactoryInstance().GetDatabase(rootCtx, "default")
	log.Printf("old db: %v,%p", db, db)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result User
		newDB := db.WithContext(rootCtx)
		//log.Printf("context:%v,%p", ctx, ctx)
		log.Printf("new db:%v,%p", newDB, newDB)
		err := newDB.GetOne(&User{}, &result, "SELECT * from #TABLE# where id =?", 1)
		assert.NoError(t, err, "")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result User
		log.Printf("context:%v,%p", rootCtx, rootCtx)
		ctx := context.WithValue(rootCtx, "ti", "test-ti-1234")
		log.Printf("context:%v,%p", ctx, ctx)
		newDB := db.WithContext(ctx)
		log.Printf("new db2:%v,%p", newDB, newDB)
		err := newDB.Create(&User{}, &result, "INSERT INTO `user` (`id`,`name`,`address`) VALUES (?,?,?)", 1, "test1", "testaddress")
		assert.NoError(t, err, "")
	}()
	wg.Wait()

}
