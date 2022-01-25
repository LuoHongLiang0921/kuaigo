// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/4/23 7:16 下午

package cache

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/file"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_cacheManager_GetCache(t *testing.T) {
	file.RegisterConfigHandler()
	configAddr := os.Getenv("CONFIG_FILE_ADDR")
	provider, err := manager.NewConfigSource(configAddr)
	if err == manager.ErrConfigAddr {
		klog.Panic(err.Error())
	}
	if err := conf.LoadFromConfigSource(provider, yaml.Unmarshal); err != nil {
		klog.Panic(err.Error())
	}
	ctx := context.Background()
	GetCacheManagerInstance().GetCache(ctx, "redis").WithContext(ctx).Set("test", 1, 10*time.Second)
	result := GetCacheManagerInstance().GetCache(ctx, "redis").WithContext(ctx).Get("test")
	assert.EqualValues(t, "1", result)

}
