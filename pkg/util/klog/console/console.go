package console

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"sync"

)

var once sync.Once

// RegisterOutputCreatorHandler
// 	@Description
func RegisterOutputCreatorHandler() {
	once.Do(func() {
		manager.Register(OutputConsole, func(cfg interface{}) []zap.Core {
			if consoleCfg, ok := cfg.(Config); ok {
				consoleCore := consoleCfg.Build()
				return []zap.Core{consoleCore}
			}
			return nil
		})
	})
}
