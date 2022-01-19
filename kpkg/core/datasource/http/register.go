package http

import (
	"context"

	"git.bbobo.com/framework/tabby/pkg/conf"
	"git.bbobo.com/framework/tabby/pkg/core/datasource/manager"
	"git.bbobo.com/framework/tabby/pkg/flag"
	"git.bbobo.com/framework/tabby/pkg/xlog"
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
			xlog.Panic(context.TODO(), "new http dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
	manager.Register(DataSourceHttp, dataSourceCreator)
	manager.Register(DataSourceHttps, dataSourceCreator)
}
