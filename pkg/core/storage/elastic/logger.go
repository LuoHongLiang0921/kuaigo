// @Description es logger 目前只能初始化客户端时设置，无法串联业务，先支持 debug 输出到本地查看问题，不记录到日志收集系统

package elastic

import (
	"fmt"
	"time"
)

var (
	ErrorLogger = Logger{
		Close: true,
		Type:  "error",
	}
	InfoLogger = Logger{
		Close: true,
		Type:  "info",
	}
	TraceLogger = Logger{
		Close: true,
		Type:  "trace",
	}
	NowFunc = func() time.Time {
		return time.Now()
	}
)

// Logger
type Logger struct {
	Close bool   // 是否关闭
	Type  string // 日志类型
}

// Printf 打印日志
// 	@Description 打印日志
// 	@receiver l Logger
//	@Param format 格式内容
//	@Param v 格式对用的值
func (l Logger) Printf(format string, v ...interface{}) {
	if !l.Close {
		fmt.Print("\n\033[33m[" + NowFunc().Format("2006-01-02 15:04:05") + "]\033[0m" +
			"\033[35m(" + l.Type + ")\033[0m  " + fmt.Sprintf(format, v...) + "\n")
	}
}
