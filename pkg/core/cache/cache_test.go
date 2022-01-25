// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/4/23 7:16 下午

package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"git.bbobo.com/framework/tabby/pkg/conf"
	"git.bbobo.com/framework/tabby/pkg/core/configsource/file"
	"git.bbobo.com/framework/tabby/pkg/core/configsource/manager"
	"git.bbobo.com/framework/tabby/pkg/util/xlog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_cacheManager_GetCache(t *testing.T) {
	file.RegisterConfigHandler()
	configAddr := os.Getenv("CONFIG_FILE_ADDR")
	provider, err := manager.NewConfigSource(configAddr)
	if err == manager.ErrConfigAddr {
		xlog.Panic(err.Error())
	}
	if err := conf.LoadFromConfigSource(provider, yaml.Unmarshal); err != nil {
		xlog.Panic(err.Error())
	}
	ctx := context.Background()
	GetCacheManagerInstance().GetCache(ctx, "redis").WithContext(ctx).Set("test", 1, 10*time.Second)
	result := GetCacheManagerInstance().GetCache(ctx, "redis").WithContext(ctx).Get("test")
	assert.EqualValues(t, "1", result)

}
