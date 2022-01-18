package ecode

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/core/net/xhttp"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	nhttp "net/http"
	"sync"

)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum int32 = 9999

const (
	// ResultMaxLimit 响应客户端，列表返回最大条数
	ResultMaxLimit = 2000
	// ResultMaxPage 最大允许页码
	ResultMaxPage = 100
	// ResultMaxPageLimit  每页充许最大条数
	ResultMaxPageLimit = 30
)

const (
	// CodeSuccess 请求成功
	CodeSuccess = 1
	//CodeCreated 资源请求已创建
	CodeCreated = 201
	//CodePartialContent 请求指定字节内容成功
	CodePartialContent = 206
	//CodeBadRequest 请求错误
	CodeBadRequest = 400
	//CodeUnauthorized  证书错误
	CodeUnauthorized = 401 //   证书错误
	//CodeForbidden  拒绝访问
	CodeForbidden = 403
	//CodeNotFound   资源不存在
	CodeNotFound = 404
	//CodeMethodNotAllowed   HTTP请求方法不支持
	CodeMethodNotAllowed = 405
	//CodeRequestedRangeNotSatisfiable   请求字节范围超出实际大小
	CodeRequestedRangeNotSatisfiable = 416
	//CodeMissingParameter  必要参数缺失  业务方可添加缺失的参数名
	CodeMissingParameter = 4001
	//CodeMissingParameterAppId   AppId缺失
	CodeMissingParameterAppId = 4002
	//CodeMissingParameterServiceId   ServiceId缺失
	CodeMissingParameterServiceId = 4003
	//CodeMissingParameterSecurityHeader   未提供认证凭据 (yanghanwei 暂时替换掉)
	//CodeMissingParameterSecurityHeader = 4004
	//CodeDataNotFound 列表接口没有返回数据
	CodeApiResponseDataNull = 4004
	//CodeInvalidParameter   不合法的参数  参数值或参数校验不合法
	CodeInvalidParameter = 4005
	//CodeInvalidParameterType   不合法的参数类型
	CodeInvalidParameterType = 4006
	//CodeInValidParameterSignature   不合法的签名
	CodeInValidParameterSignature = 4007
	//CodeInvalidParameterSecurity   认证实效，请重新登录
	CodeInvalidParameterSecurity = 4008
	//CodeInvalidParameterUser   不合法的用户
	CodeInvalidParameterUser = 4009
	//CodeInvalidParameterTimestamp   不合法的时间戳
	CodeInvalidParameterTimestamp = 4010
	//CodeInvalidParameterAppId   不合法的产品id    产品id不存在，id长度不合法
	CodeInvalidParameterAppId = 4011
	//CodeInvalidParameterServiceId   不合法的服务id    id长度不合法
	CodeInvalidParameterServiceId = 4012
	//CodeInvalidParameterPhone   不合法的手机号
	CodeInvalidParameterPhone = 4013
	//CodeInvalidParameterFormat   不合法的格式
	CodeInvalidParameterFormat = 4014
	//CodeInvalidVerificationCode   验证码失效
	CodeInvalidVerificationCode = 4015
	//CodeInvalidDevice   不合法的设备
	CodeInvalidDevice = 4016
	//CodeRequestRateLimit   请求频次超出限制
	CodeRequestRateLimit = 4017
	//CodeIPLimit   IP访问限制
	CodeIPLimit = 4018
	//CodeApiRequestLimit   应用Api请求限制   接口调用过于频繁
	CodeApiRequestLimit = 4019
	//CodeUserLimit  用户限制
	CodeUserLimit = 4020
	//CodePhoneLimit  手机号限制
	CodePhoneLimit = 4021
	//CodePermissionDenied   拒绝访问    签名、鉴权通过，但权限不足
	CodePermissionDenied = 4022
	//CodeUserAccountNotExist   用户账号不存在
	CodeUserAccountNotExist = 4023
	//CodeUserNotLogin   用户未登录
	CodeUserNotLogin = 4024
	//CodeUserNotSetPassword   用户密码未设置
	CodeUserNotSetPassword = 4025
	//CodePasswordError   密码错误
	CodePasswordError = 4026
	//CodePasswordsNotConsistent   两次密码不一致
	CodePasswordsNotConsistent = 4027
	//CodeOldPasswordError   旧密码输入错误
	CodeOldPasswordError = 4028
	//CodeEmptyPassword   密码为空
	CodeEmptyPassword = 4029
	//CodeRequestNextIdError   获取新id错误
	CodeRequestNextIdError = 4030
	//CodeAuthTokenTimeout   访问令牌过期
	CodeAuthTokenTimeout = 4031
	//CodeInvalidAuthToken   无效的访问令牌
	CodeInvalidAuthToken = 4032

	// CodeHeaderPParamError p参数不合法
	CodeHeaderPParamError = 4033
	// CodeHeaderPParamError s参数不合法
	CodeHeaderSParamError = 4034
	//CodeInternalServerError   服务器错误
	CodeInternalServerError = 500
	//CodeServiceUnavailable   服务暂时不可用
	CodeServiceUnavailable = 5001
	//CodeRemoteServerError   远程服务错误
	CodeRemoteServerError = 5002
	//CodeServerBusy   服务器繁忙
	CodeServerBusy = 5003
	//CodeRequestTimeout   请求处理超时
	CodeRequestTimeout = 5004
	//CodeUnknownError   服务器未知错误
	CodeUnknownError = 5005
	//CodeDatabaseError   数据库错误
	CodeDatabaseError = 5006
	//CodeDatabaseCreateError   数据库信息创建错误
	CodeDatabaseCreateError = 5007
	//CodeDatabaseUpdateError   数据库信息更新错误
	CodeDatabaseUpdateError = 5008
	//CodeDatabaseQueryError   数据库信息查询错误
	CodeDatabaseQueryError = 5009
	//CodeDatabaseDeleteError   数据库信息删除错误
	CodeDatabaseDeleteError = 5010
	//CodeCacheError   缓存处理错误
	CodeCacheError = 5011
	//CodeFlushError   缓冲区操作错误
	CodeFlushError = 5012
	//CodeMQPushError   推送MQ信息错误
	CodeMQPushError = 5013
)

var (
	codeRW        sync.RWMutex
	codeMsgMapper = make(map[int]string)
)

type Config struct {
	//工具箱中状态码服务地址
	StatusCodeURL string
}

func (c *Config) Build(key string) *ECode {
	url := conf.GetString(key)
	return &ECode{
		StatusCodeURL: url,
	}
}

type ECode struct {
	//CodeMsgMapper code,msg 相互装换
	CodeMsgMapper map[int]string
	//appId  应用服务ID
	AppId string `json:"app_id"`
	//BusinessID 业务id
	BusinessID string `json:"businessId"`
	//ServiceName 发送日志服务名称
	ServiceName string `json:"service_name"`
	//工具箱中状态码服务地址
	StatusCodeURL string `json:"statusCodeUrl"`
	//工具箱中状态码更新地址
	StatusCodeUpdateUrl string `json:"statusCodeUpdateUrl"`
	//工具箱中状态码创建地址
	StatusCodeCreateUrl string `json:"statusCodeCreateUrl"`
	//状态码列表查询地址
	StatusCodeListUrl string `json:"statusCodeListUrl"`
}

type statusCodeListReq struct {
	ServiceId int `json:"serviceId"`
}

type statusCodeReq struct {
	Code int `json:"code"`
}

type statusCreateCodeReq struct {
	AppId     int    `json:"appId"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Lang      string `json:"lang"`
	Name      string `json:"name"`
	ServiceId int    `json:"serviceId"`
}

type statusCodeListResp struct {
	Code   int          `json:"code"`
	Msg    string       `json:"msg"`
	Expire int64        `json:"expire"`
	Data   []statusCode `json:"data"`
	RT     int64        `json:"rt"`
}

type statusCodeResp struct {
	Code   int        `json:"code"`
	Msg    string     `json:"msg"`
	Expire int64      `json:"expire"`
	Data   statusCode `json:"data"`
	RT     int64      `json:"rt"`
}

type statusCode struct {
	Code       int    `json:"code"`
	ServiceId  int    `json:"serviceId"`
	Name       string `json:"name"`
	Lang       string `json:"lang"`
	Msg        string `json:"msg"`
	CreatTime  int64  `json:"creatTime"`
	UpdateTime int64  `json:"updateTime"`
}

func (c *ECode) LoadCreateCodeMsg(ctx context.Context, appId int, code int, msg string, lang string, name string, serviceId int) {
	header := nhttp.Header{}
	header.Set("Content-Type", "application/json")
	req := &statusCreateCodeReq{
		AppId:     appId,
		Code:      code,
		Msg:       msg,
		Lang:      lang,
		Name:      name,
		ServiceId: serviceId,
	}
	resp := &statusCodeResp{}
	err := xhttp.PostWithUnmarshal(ctx, nil, c.StatusCodeCreateUrl, header, req, resp)
	if err != nil {
		klog.Errorf("", klog.FieldCommon(nil), klog.FieldParams(map[string]interface{}{
			"msg": "LoadCreateCodeMsg",
		}), klog.String("type", "bizLog"))
		return
	}
	if c.CodeMsgMapper == nil {
		c.CodeMsgMapper = make(map[int]string)
	}
	c.CodeMsgMapper[resp.Data.Code] = resp.Data.Msg
}

func (c *ECode) LoadUpdateCodeMsg(ctx context.Context, appId int, code int, msg string, lang string, name string, serviceId int) {
	header := nhttp.Header{}
	header.Set("Content-Type", "application/json")
	req := &statusCreateCodeReq{
		AppId:     appId,
		Code:      code,
		Msg:       msg,
		Lang:      lang,
		Name:      name,
		ServiceId: serviceId,
	}
	resp := &statusCodeResp{}
	err := xhttp.PostWithUnmarshal(ctx, nil, c.StatusCodeUpdateUrl, header, req, resp)
	if err != nil {
		klog.Errorf("", klog.FieldCommon(nil), klog.FieldParams(map[string]interface{}{
			"msg": "LoadUpdateCodeMsg",
		}), klog.String("type", "bizLog"))
		return
	}
	if c.CodeMsgMapper == nil {
		c.CodeMsgMapper = make(map[int]string)
	}
	c.CodeMsgMapper[resp.Data.Code] = resp.Data.Msg
}

// CreateCode 状态码创建
func (c *ECode) LoadCodeMsg(ctx context.Context, code int) {
	header := nhttp.Header{}
	header.Set("Content-Type", "application/json")
	req := &statusCodeReq{
		Code: code,
	}
	resp := &statusCodeResp{}
	err := xhttp.PostWithUnmarshal(ctx, nil, c.StatusCodeURL, header, req, resp)
	if err != nil {
		klog.Errorf("", klog.FieldCommon(nil), klog.FieldParams(map[string]interface{}{
			"msg": "LoadCodeMsg",
		}), klog.String("type", "bizLog"))
		return
	}
	if c.CodeMsgMapper == nil {
		c.CodeMsgMapper = make(map[int]string)
	}
	c.CodeMsgMapper[resp.Data.Code] = resp.Data.Msg
}

//LoadCodeListMsg 状态码列表服务
func (c *ECode) LoadCodeListMsg(ctx context.Context, serviceID int) {
	header := nhttp.Header{}
	header.Set("Content-Type", "application/json")
	req := &statusCodeListReq{
		ServiceId: serviceID,
	}
	var resp statusCodeListResp
	err := xhttp.PostWithUnmarshal(ctx, nil, c.StatusCodeListUrl, header, req, &resp)
	if err != nil {
		klog.Errorf("", klog.FieldCommon(nil), klog.FieldParams(map[string]interface{}{
			"msg": "LoadCodeMsg",
		}), klog.String("type", "bizLog"))
		return
	}
	codeRW.Lock()
	defer codeRW.Unlock()
	if c.CodeMsgMapper == nil {
		c.CodeMsgMapper = make(map[int]string)
	}
	for _, sc := range resp.Data {
		c.CodeMsgMapper[sc.Code] = sc.Msg
	}
}

// GetMsg ...
// 获取code 对应的msg
func (c *ECode) GetMsg(ctx context.Context, code int) string {
	codeRW.RLock()
	defer codeRW.RUnlock()
	if v, ok := c.CodeMsgMapper[code]; ok {
		return v
	}
	return ""
}
