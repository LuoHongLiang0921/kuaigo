package khttp

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kjson"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type (
	// Client R
	RClient = resty.Client
	// User
	User = resty.User
	//Response = resty.Response
	Request = resty.Request
)
type (
	Response struct {
		*resty.Response
	}
)
type (
	RetryConditionFunc = resty.RetryConditionFunc
	RequestMiddleware  = resty.RequestMiddleware
	ResponseMiddleware = resty.ResponseMiddleware
	RequestLogCallback = resty.RequestLogCallback
	RequestLog         = resty.RequestLog
	ResponseLog        = resty.ResponseLog
)

var (
	New              = resty.New
	NewWithClient    = resty.NewWithClient
	NewWithLocalAddr = resty.NewWithLocalAddr
	RetryConditions  = resty.RetryConditions
)

// retryCondition 重试条件
// 	@Description:  当出错的时候重试
// 	@return RetryConditionFunc
func retryCondition() RetryConditionFunc {
	return func(response *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return false
	}
}

//Json
// @Description:响应内容转json
// @Receiver resp
// @Param v
// @Return error
func (resp *Response) Json(v interface{}) error {
	if resp.String() == "" {
		resp.String()
	}
	return kjson.DecodeFromString(resp.String(), v)

}

// IsUnauthorized
//  @Description: 返回是否是认证页面
//  @Receiver resp
//  @Return bool
func (resp *Response) IsUnauthorized() bool {
	return resp.StatusCode() == http.StatusUnauthorized
}
