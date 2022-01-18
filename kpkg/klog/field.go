package klog

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (

	//LogTypeRunning 运行日志
	LogTypeRunning = "running"
	//LogTypeError 错误日志
	LogTypeError = "error"
	//LogTypeAccess 访问日志
	LogTypeAccess = "access"
	//LogTypeTask 任务日志
	LogTypeTask = "task"

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

type Log struct {
	//common放通用的头部参数
	Common Common `json:"common"`
	//params放具体的业务数据
	Params map[string]interface{} `json:"params"`
	//LogType type表示什么业务的日志
	LogType string `json:"type"`
}

// 应用唯一标识符
func FieldAid(value string) Field {
	return String("aid", value)
}

// 模块
func FieldMod(value string) Field {
	value = strings.Replace(value, " ", ".", -1)
	return String("mod", value)
}

// 依赖的实例名称。以mysql为例，"dsn = "root:juno@tcp(127.0.0.1:3306)/juno?charset=utf8"，addr为 "127.0.0.1:3306"
func FieldAddr(value string) Field {
	return String("addr", value)
}

// FieldAddrAny ...
func FieldAddrAny(value interface{}) Field {
	return Any("addr", value)
}

// FieldName ...
func FieldName(value string) Field {
	return String("name", value)
}

// FieldType ...
func FieldType(value string) Field {
	return String("type", value)
}

// FieldCode ...
func FieldCode(value int32) Field {
	return Int32("code", value)
}

// 耗时时间
func FieldCost(value time.Duration) Field {
	return String("cost", fmt.Sprintf("%.3f", float64(value.Round(time.Microsecond))/float64(time.Millisecond)))
}

// FieldKey ...
func FieldKey(value string) Field {
	return String("key", value)
}

// 耗时时间
func FieldKeyAny(value interface{}) Field {
	return Any("key", value)
}

// FieldValue ...
func FieldValue(value string) Field {
	return String("value", value)
}

// FieldValueAny ...
func FieldValueAny(value interface{}) Field {
	return Any("value", value)
}

// FieldErrKind ...
func FieldErrKind(value string) Field {
	return String("errKind", value)
}

// FieldErr ...
func FieldErr(err error) Field {
	return zap.Error(err)
}

// FieldErr ...
func FieldStringErr(err string) Field {
	return String("err", err)
}

// FieldExtMessage ...
func FieldExtMessage(vals ...interface{}) Field {
	return zap.Any("ext", vals)
}

// FieldStack ...
func FieldStack(value []byte) Field {
	return ByteString("stack", value)
}

// FieldMethod ...
func FieldMethod(value string) Field {
	return String("method", value)
}

// FieldEvent ...
func FieldEvent(value string) Field {
	return String("event", value)
}

//FieldParams ...
func FieldParams(value interface{}) Field {
	return Any(constant.ParamsKey, value)
}

//FieldCommon ...
func FieldCommon(value interface{}) Field {
	return Any(constant.CommonKey, value)
}

//FiledLogType ...
func FiledLogType(value string) Field {
	return String(constant.TypeKey, value)
}
