// @Description 可选配置项

package khttp

import (
	"net/http"
)

// RequestOption 请求选项
type RequestOption struct {
	hc                *http.Client      // 标准库client
	header            http.Header       // 请求头
	timeOut           int               // 超时处理
	retryCount        int               // 重试次数
	isReturnRaw       bool              // 原始内容返回，不处理，默认需要处理原始内容
	isRequestFormData bool              // 请求内容是否为formData 默认不是form data
	queryParams       map[string]string // 查询参数
	pathParams        map[string]string // 路径参数
	getReqBody        interface{}       //查询方法请求体
}

// Option
type Option func(r *RequestOption)

// WithHttpClient
// 	@Description 设置标准库http client
//	@Param client 标准库http client
// 	@Return *Option 设置标准库 http client 后的 Option 函数
func WithHttpClient(client *http.Client) Option {
	return func(r *RequestOption) {
		r.hc = client
	}
}

// WithTimeOut
// 	@Description 超时设置
//	@Param timeout 超时 单位为毫秒
// 	@Return *Option 设置超时后的 Option
func WithTimeOut(timeout int) Option {
	return func(r *RequestOption) {
		r.timeOut = timeout
	}
}

// WithRetryCount
// 	@Description  设置重试次数
//	@Param retryCount 重试次数
// 	@Return *Option 设置重试次数后的 Option
func WithRetryCount(retryCount int) Option {
	return func(r *RequestOption) {
		r.retryCount = retryCount
	}
}

// WithHttpHeader
// 	@Description 设置请求头
//	@Param header 请求头
// 	@Return *Option 设置请求头后的 Option
func WithHttpHeader(header http.Header) Option {
	return func(r *RequestOption) {
		r.header = header
	}
}

// WithIsReturnRaw
// 	@Description 设置是否返回原始内容
//	@Param isReturnRaw true 返回原始内容，false
// 	@Return Option 设置否返回原始内容的 Option
func WithIsReturnRaw(isRaw bool) Option {
	return func(r *RequestOption) {
		r.isReturnRaw = isRaw
	}
}

// WithIsFormData
// 	@Description  设置是否是form 表单请求内容
//	@Param isFormData
// 	@Return Option 设置是否是form 表单请求内容的 Option
func WithIsFormData(isFormData bool) Option {
	return func(r *RequestOption) {
		r.isRequestFormData = isFormData
	}
}

// WithQueryParams
// 	@Description 设置查询参数
//	@Param queryParams
// 	@Return Option
func WithQueryParams(queryParams map[string]string) Option {
	return func(r *RequestOption) {
		r.queryParams = queryParams
	}
}

// WithPathParams
// 	@Description 设置路径参数
//	@Param pathParams
// 	@Return Option
func WithPathParams(pathParams map[string]string) Option {
	return func(r *RequestOption) {
		r.pathParams = pathParams
	}
}

// WithGetReqBody
// 	@Description 设置 http get 请求方法中请求内容
//	@Param b
// 	@Return Option
func WithGetReqBody(b interface{}) Option {
	return func(r *RequestOption) {
		r.getReqBody = b
	}
}

// Build
// 	@Description 构造请求选项
//	@Param opts 请求参数
// 	@Return *RequestOption 构造后的请求选项
func Build(opts ...Option) *RequestOption {
	var o RequestOption
	for _, f := range opts {
		f(&o)
	}
	return &o
}
