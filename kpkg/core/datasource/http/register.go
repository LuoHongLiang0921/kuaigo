package http

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/datasource/manager"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
)

// Defines http/https scheme
const (
	DataSourceHttp  = "http"
	DataSourceHttps = "https"
)

func init() {
	dataSourceCreator := func() conf.DataSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			klog.Panic(context.TODO(), "new http dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
	manager.Register(DataSourceHttp, dataSourceCreator)
	manager.Register(DataSourceHttps, dataSourceCreator)
}
