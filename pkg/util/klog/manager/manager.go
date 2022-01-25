package manager

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
)

// OutputCreatorFunc 创建输出源实例函数
type OutputCreatorFunc func(cfg interface{}) []zap.Core

var (
	registry = make(map[string]OutputCreatorFunc)
)

// Register
// 	@Description  注册 output 类型的 输出源
//	@Param output
//	@Param creator
func Register(output string, creator OutputCreatorFunc) {
	registry[output] = creator
}

// GetCreator
// 	@Description 获取 output 类型的 创建函数
//	@Param output
// 	@Return bool
// 	@Return OutputCreatorFunc
func GetCreator(output string) (bool, OutputCreatorFunc) {
	v, ok := registry[output]
	return ok, v
}
