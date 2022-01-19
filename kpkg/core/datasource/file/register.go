package file

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/datasource/manager"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func init() {
	manager.Register(DataSourceFile, func() conf.DataSource {
		var (

			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			klog.Panic(context.TODO(), "new file dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	})
	manager.DefaultScheme = DataSourceFile
}
