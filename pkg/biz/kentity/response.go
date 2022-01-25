package kentity

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/kgin"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
	"net/http"

)

const (
	CodeOk = 1
	MsgOk  = "SUCCESS"
)

//ResponseOption ...
type ResponseOption func(response *Response)

//Response ...
// see: http://wiki.yixiahd.com/pages/viewpage.action?pageId=5734480
type Response struct {
	// Code 请求结果
	Code int `json:"code"`
	//Msg 结果描述串
	Msg string `json:"msg"`
	//Data 返回数据
	Data interface{} `json:"data"`
	//扩展数据
	Ext interface{} `json:"ext,omitempty"`
	//服务端响应时间
	RT int `json:"rt,omitempty"`
	//数据过期时间 毫秒
	Expire int `json:"expire,omitempty"`
	//服务端data部分数据签名
	SK string `json:"sk,omitempty"`
}

//ListData 分页数据
type ListData struct {
	//Page  当前数据所在页
	Page int `json:"page"`
	//Limit 每页显示条数
	Limit int `json:"limit"`
	//Count 本次返回数据条数
	Count int `json:"count"`
	//Cursor 游标值
	Cursor int64 `json:"cursor"`
	//Total 服务器估算数据总条数
	Total int `json:"total"`
	//List 数组列表
	List interface{} `json:"list"`
}

func WithCode(code int, msg string) ResponseOption {
	return func(response *Response) {
		response.Code = code
		response.Msg = msg
	}
}

func WithExt(ext interface{}) ResponseOption {
	return func(response *Response) {
		response.Ext = ext
	}
}

func WithRT(rt int) ResponseOption {
	return func(response *Response) {
		response.RT = rt
	}
}

func WithSK(sk string) ResponseOption {
	return func(response *Response) {
		response.SK = sk
	}
}

func WithExpire(expire int) ResponseOption {
	return func(response *Response) {
		response.Expire = expire
	}
}

//NewResponse ...
func NewResponse(data interface{}, opts ...ResponseOption) *Response {
	r := new(Response)
	for _, f := range opts {
		f(r)
	}
	r.Data = data
	return r
}

func JSONResponse(c *kgin.TContext, resp *Response) {
	c.Set(constant.RespCode, resp.Code)
	c.JSON(http.StatusOK, &resp)
}

// @Description:
// @Param c
// @Param resp
func AbortWithJSONResponse(c *kgin.TContext, resp *Response) {
	c.Set(constant.RespCode, resp.Code)
	c.Abort()
	c.JSON(http.StatusOK, &resp)
}

// JSON
//  @Description 正常返回响应
//  @Param c
//  @Param data
func JSON(c *kgin.TContext, data interface{}) {
	var resp Response
	resp.Code = CodeOk
	resp.Data = data
	resp.Msg = MsgOk
	c.Set(constant.RespCode, CodeOk)
	c.JSON(http.StatusOK, &resp)
}

// ErrJSON
//  @Description 根据错误码返回错误
//  @Param c
//  @Param code
func ErrJSON(c *kgin.TContext, code int) {
	var resp Response
	resp.Code = code
	//msg := xconfig.GetMsg(code)
	resp.Msg = ""
	c.Set(constant.RespCode, code)
	c.JSON(http.StatusOK, &resp)
}

// ErrorJSON TODO 重命名
//  @Description 自定义err返回错误
//  @Param c
//  @Param err
func ErrorJSON(c *kgin.TContext, err error) {
	var resp Response
	if v, ok := err.(errs.Error); ok {
		resp.Code = v.Code
		c.Set(constant.RespCode, v.Code)
		resp.Msg = v.Msg
	}
	//todo: 从配置服务获取
	//msg := xconfig.GetMsg(code)
	//resp.Msg = ""
	c.JSON(http.StatusOK, &resp)
}

// AbortWithErrorJSON
// @Description:
// @Param c gin context
// @Param err err
func AbortWithErrorJSON(c *kgin.TContext, err error) {
	var resp Response
	if v, ok := err.(errs.Error); ok {
		resp.Code = v.Code
		c.Set(constant.RespCode, v.Code)
		resp.Msg = v.Msg
	}
	//todo: 从配置服务获取
	//msg := xconfig.GetMsg(code)
	//resp.Msg = ""
	c.Abort()
	c.JSON(http.StatusOK, &resp)
}
