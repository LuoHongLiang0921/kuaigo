package kdebug

import (
	"flag"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kstring"
	"os"
	"sync"

	"github.com/tidwall/pretty"
)

var (
	isTestingMode     bool
	isDevelopmentMode = os.Getenv(constant.TabbyMode) == constant.TabbyModeDev
)

func init() {
	if isDevelopmentMode {
		klog.RunningLogger.SetLevel(klog.DebugLevel)
		klog.KuaigoLogger.SetLevel(klog.DebugLevel)
	}
}

// IsTestingMode 判断是否在测试模式下
var onceTest = sync.Once{}

// IsTestingMode ...
func IsTestingMode() bool {
	onceTest.Do(func() {
		isTestingMode = flag.Lookup("test.v") != nil
	})

	return isTestingMode
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return isDevelopmentMode || isTestingMode
}

// IfPanic ...
func IfPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyJsonPrint ...
func PrettyJsonPrint(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		kcolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(kstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}
