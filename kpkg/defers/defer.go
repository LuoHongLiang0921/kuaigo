package defers

import (
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kdefer"
)

var (
	globalDefers = kdefer.NewStack()
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Clean 清除
func Clean() {
	globalDefers.Clean()
}
