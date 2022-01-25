// @Description

package klog

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"regexp"
	"runtime"
	"strings"

	"github.com/pborman/uuid"
)

var logSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	logSourceDir = regexp.MustCompile(`log_func\.go`).ReplaceAllString(file, "")
}

// makeFields
// 	@Description  根据上下文构造日志字段数组
// 	@Receiver logger Logger
//	@Param ctx 上下文
//	@Param level 日志级别
// 	@Return []Field 构造后的日志字段数组
func (lg *Logger) makeFields(level zap.Level) []zap.Field {
	if !lg.isPrintCommon {
		return []zap.Field{}
	}

	ctx := lg.getContext()
	com, ok := FromContext(ctx)
	var fields []zap.Field
	file, line := fileWithLineNum()
	if !ok {
		com = Common{}
	}
	traceId := com.TraceId
	if traceId == "" {
		traceId = uuid.New()
	}
	serviceName := lg.config.ServiceName
	if serviceName == "" {
		serviceName = com.ServiceName
	}
	fields = append(fields, zap.Int("appId", com.AppId))
	fields = append(fields, zap.String("traceId", traceId))
	fields = append(fields, zap.String("serviceSource", com.ServiceSource))
	fields = append(fields, zap.String("serviceName", serviceName))
	switch lg.loggerType {
	case LogTypeError:
		fallthrough
	case LogTypeRunning:
		fields = append(fields, zap.String("fileName", file))
		fields = append(fields, zap.String("logLevel", level.String()))
		fields = append(fields, zap.Int("line", line))
		fields = append(fields, zap.String("requestIp", com.RequestIp))
		fields = append(fields, zap.String("requestUri", com.RequestUri))
	case LogTypeAccess:
		fields = append(fields, zap.String("requestIp", com.RequestIp))
		fields = append(fields, zap.String("requestUri", com.RequestUri))
		fields = append(fields, zap.Int("processCode", com.ProcessCode))
		fields = append(fields, zap.Int("costTime", int(com.CostTime)))
		fields = append(fields, zap.Int("code", com.Code))
		fields = append(fields, zap.Any("p", com.P))
		fields = append(fields, zap.String("uid", com.UID))
		fields = append(fields, zap.String("msg", ""))
		fields = append(fields, zap.String("logLevel", ""))
	case LogTypeTask:
		fields = append(fields, zap.String("fileName", file))
		fields = append(fields, zap.String("logLevel", level.String()))
		fields = append(fields, zap.Int("line", line))
	}
	return fields
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

func fileWithLineNum() (string, int) {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, logSourceDir) || strings.HasSuffix(file, "_test.go")) {
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file, line
			}
			idx = strings.LastIndexByte(file[:idx], '/')
			if idx == -1 {
				return file, line
			}
			return file[idx+1:], line
		}
	}
	return "", 0
}
