// @Description

package file

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"sync"
)

var once sync.Once

// RegisterOutputCreatorHandler
// 	@Description 注册文件
//	@Param c
func RegisterOutputCreatorHandler() {
	once.Do(func() {
		manager.Register(OutputFile, func(cfg interface{}) []zap.Core {
			if fileCfg, ok := cfg.(Config); ok {
				fileCore := fileCfg.Build()
				return []zap.Core{fileCore}
			}
			return nil
		})
	})
}
