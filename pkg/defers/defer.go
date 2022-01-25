// @Description

package defers

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kdefer"
)

var (
	globalDefers = kdefer.NewStack()
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Execute
// 	@Description 执行栈中函数
func Execute() {
	globalDefers.Execute()
}
