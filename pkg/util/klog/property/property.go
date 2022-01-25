package property

import (
	"sort"
	"strings"
	"sync"
)

var (
	// 内置字段
	buildInFieldsMap = map[string]string{
		"appId":         "ai",
		"traceId":       "ti",
		"serviceSource": "ss",
		"serviceName":   "sn",
		"fileName":      "fn",
		"line":          "l",
		"logLevel":      "lv",
		"requestIp":     "i",
		"requestUri":    "u",
		"msg":           "msg",
		"processCode":   "pc",
		"costTime":      "ct",
		"code":          "cd",
		"p":             "p",
	}
	reverseFieldsMap = map[string]string{
		"ai":  "appId",
		"ti":  "traceId",
		"ss":  "serviceSource",
		"sn":  "serviceName",
		"fn":  "fileName",
		"l":   "line",
		"lv":  "logLevel",
		"i":   "requestIp",
		"u":   "requestUri",
		"msg": "msg",
		"pc":  "processCode",
		"ct":  "costTime",
		"cd":  "code",
		"p":   "p",
	}
	mu sync.RWMutex
)

type Property struct {
	K string
	v interface{}
}

// RegisterProperty 注册新的属性变量
// 	@Description
//	@Param k
//	@Param v
func RegisterProperty(k string, v string) {
	mu.Lock()
	buildInFieldsMap[k] = v
	mu.Unlock()
}

// UnRegisterProperty 删除新的属性变量
// 	@Description
//	@Param k
func UnRegisterProperty(k string) {
	mu.Lock()
	delete(buildInFieldsMap, k)
	mu.Unlock()
}

// BuildInExpress
// 	@Description 获取 内置字段 表达式字符串值
//	@Param sored
// 	@Return string
func BuildInExpress(sored bool) string {
	var builder strings.Builder
	if !sored {
		for k := range buildInFieldsMap {
			builder.WriteString(k + "=" + "%" + "{" + k + "}" + " ")
		}
		return builder.String()
	}

	keys := make([]string, 0, len(buildInFieldsMap))
	for k := range buildInFieldsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		builder.WriteString(v + "=" + "%" + "{" + v + "}" + " ")
	}
	return builder.String()
}
