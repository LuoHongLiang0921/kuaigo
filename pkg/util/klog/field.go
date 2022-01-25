// @Description 日志字段

package klog

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"strings"
	"time"
)

const (
	LogKeyRunning = "running"
	LogKeyError   = "error"
	LogKeyAccess  = "access"
	LogKeyTask    = "task"

	//LogTypeRunning 运行日志
	LogTypeRunning = "running"
	//LogTypeError 错误日志
	LogTypeError = "error"
	//LogTypeAccess 访问日志
	LogTypeAccess = "access"
	//LogTypeTask 任务日志
	LogTypeTask = "task"
	// LogTypeTabby 框架日志
	LogTypeTabby = "tabby"

	//LogTypeBigData BI、推荐和广告上报日志
	LogTypeBigData = "bigdataLog"
	//ProcessCodeRequest 接受请求日志
	ProcessCodeRequest = 0
	//ProcessCodeResponse 响应返回日志
	ProcessCodeResponse = 1

	//ServiceSourceApp 日志来源为app
	ServiceSourceApp = "app"
	//ServiceSourceWeb 日志来源为web
	ServiceSourceWeb = "web"
	//ServiceSourceH5 日志来源为h5
	ServiceSourceH5 = "h5"
	//ServiceSourceAdmin 日志来源为admin
	ServiceSourceAdmin = "admin"
	// ServiceSourceTask 日志来源为task
	ServiceSourceTask = "task"
	// ServiceSourceConsumer 日志来源为 consumer
	ServiceSourceConsumer = "consumer"
	// ServiceSourceService 日志来源为 service
	ServiceSourceService = "service"
)

//Common ...
//see:http://wiki.yixiahd.com/pages/viewpage.action?pageId=5734565
type Common struct {
	//appId  应用服务ID
	AppId int `json:"appId"`
	//TraceId 追踪id，客户端每次发起请求，需要生成一个随机不重复的Id
	TraceId string `json:"traceId"`
	//LogLevel 日志级别，有效值：DEBUG，INFO，WARN，ERROR，FATAL
	LogLevel string `json:"logLevel"`
	//ServiceSource 日志来源，有效值：app，web，h5，admin，task，consumer
	ServiceSource string `json:"serviceSource"`
	//ServiceName 发送日志服务名称
	ServiceName string `json:"serviceName"`
	//FileName 日志所在代码文件名
	FileName string `json:"fileName"`
	//Line 日志所在代码的行数
	Line int `json:"line"`
	//RequestIp 请求ip
	RequestIp string `json:"requestIp"`
	//RequestUri 请求URI
	RequestUri string `json:"requestUri"`
	//Timestamp 13位时间戳
	Timestamp int64 `json:"timestamp"`
	//ProcessCode accessLog时有效，0 代表接受请求日志，1 代表响应返回日志
	ProcessCode int `json:"processCode,omitempty"`
	//CostTime processCode=1时有效，单位毫秒
	CostTime int64 `json:"costTime,omitempty"`
	//Code processCode=1时有效，状态码，详见单独定义错误码定义字典
	Code int `json:"code,omitempty"`
	//UID 用户id
	UID string `json:"uid,omitempty"`
	//p参数，json格式（由上报服务脱敏——去掉ak在 serviceSource = service 时（RPC微服务），p传空
	P interface{} `json:"p,omitempty"`
}

// Log 日志
type Log struct {
	//common放通用的头部参数
	Common Common `json:"common"`
	//params放具体的业务数据
	Params map[string]interface{} `json:"params"`
	//LogType type表示什么业务的日志
	LogType string `json:"type"`
}

// FieldMod
// 	@Description 模块
//	@Param value 模块值
// 	@Return Field 设置模块后的字段
func FieldMod(value string) zap.Field {
	value = strings.Replace(value, " ", ".", -1)
	return zap.String("mod", value)
}

//
// FieldAddr
// 	@Description 设置地址字段
//	 依赖的实例名称。以mysql为例，"dsn = "root:juno@tcp(127.0.0.1:3306)/juno?charset=utf8"，addr为 "127.0.0.1:3306"
//	@Param value 地址值
// 	@Return Field 设置地址后的字段
func FieldAddr(value string) zap.Field {
	return zap.String("addr", value)
}

// FieldName
// 	@Description 设置 name 字段值
//	@Param value name 值
// 	@Return Field 设置name后的字段
func FieldName(value string) zap.Field {
	return zap.String("name", value)
}

// FieldCost
// 	@Description 设置耗时时间字段值
//	@Param value 耗时时间
// 	@Return Field 设置耗时后的字段
func FieldCost(value time.Duration) zap.Field {
	return zap.String("cost", fmt.Sprintf("%.3f", float64(value.Round(time.Microsecond))/float64(time.Millisecond)))
}

// FieldKey
// 	@Description 设置 key 字段值
//	@Param value key 值
// 	@Return Field 设置key字段值后的字段
func FieldKey(value string) zap.Field {
	return zap.String("key", value)
}

// FieldValueAny
// 	@Description 设置 value 字段值
//	@Param value value 值
// 	@Return Field 设置 value 字段值后的字段
func FieldValueAny(value interface{}) zap.Field {
	return zap.Any("value", value)
}

// FieldErrKind
// 	@Description 设置 errKind 字段值
//	@Param value errKind 字段值
// 	@Return Field 设置 errKind 字段后的字段
func FieldErrKind(value string) zap.Field {
	return zap.String("errKind", value)
}

// FieldErr
// 	@Description 设置error 字段值
//	@Param err error 值
// 	@Return Field 设置 error 后的字段
func FieldErr(err error) zap.Field {
	return zap.Error(err)
}

// FieldExtMessage
// 	@Description 设置 ext 字段值
//	@Param vals 要设置的值数组
// 	@Return Field 设置 ext 字段值
func FieldExtMessage(vals ...interface{}) zap.Field {
	return zap.Any("ext", vals)
}

// FieldMethod
// 	@Description 设置 method 字段值
//	@Param value method 值
// 	@Return Field 设置 method 值后的字段
func FieldMethod(value string) zap.Field {
	return zap.String("method", value)
}

// FieldStack
// 	@Description 设置stack 字段
//	@Param value stack 值
// 	@Return Field 设置stack 后的字段
func FieldStack(value []byte) zap.Field {
	return zap.ByteString("stack", value)
}

// FieldEvent
// 	@Description 设置 event 字段值
//	@Param value event 值
// 	@Return Field 设置 event 值后的字段
func FieldEvent(value string) zap.Field {
	return zap.String("event", value)
}

// FieldParams
// 	@Description 设置params 字段值
//	@Param value params 值
// 	@Return Field 设置params 后的字段
func FieldParams(value interface{}) zap.Field {
	return zap.Any(constant.ParamsKey, value)
}

// FieldCommon
// 	@Description 设置common 字段值
//	@Param value common 字段值
// 	@Return Field 设置common后字段
func FieldCommon(value interface{}) zap.Field {
	return zap.Any(constant.CommonKey, value)
}
