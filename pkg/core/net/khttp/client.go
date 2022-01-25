// @Description http请求封装库

package khttp

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"net/http"
	"time"
)

const (
	plainTextType   = "text/plain; charset=utf-8"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

var (
	// 默认超时时间 500 毫秒
	defaultTimeout = 2000
	//默认ua
	defaultUserAgent = "Tabby Request sndks/" + pkg.GetTabbyVersion()
	// 默认重试次数
	defaultRetryCount = 2
	// 如果客户端只需要访问一个host，那么最好将MaxIdleConnsPerHost与MaxIdleConns设置为相同，这样逻辑更加清晰
	//defaultMaxIdleConns = 100
	// 每个host的idle状态的最大连接数目，即idleConn中的key对应的连接数
	defaultMaxIdleConnsPerHost = 2048
	//连接保持idle状态的最大时间，超时关闭pconn
	defaultIdleConnTimeout = 90 * time.Second
)

// Client ...
type Client struct {
	Debug               bool
	TimeOut             time.Duration
	Ua                  string
	RetryCount          int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
	// IsBiz true 为业务使用,false 不是使用
	IsBiz   bool
	IsTrace bool
	httpReq *Request
	ctx     context.Context
	Cookies []*http.Cookie
}

type Header map[string]string
type Cookie map[string]string
type Params map[string]string
type Datas map[string]string // for post form

type Files map[string]string // demo:name ,filename

type Auth []string // demo:{username,password}
type Token string  // demo: BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F

func Requests(ctx context.Context) *Client {
	return &Client{
		Debug:      false,
		Ua:         defaultUserAgent,
		TimeOut:    time.Duration(defaultTimeout),
		RetryCount: defaultRetryCount,
		IsBiz:      true,
		IsTrace:    true,
		ctx:        ctx,
	}
}

// createClient
//  @Description: 创建client
//  @Receiver h
//  @Param ctx
//  @Return *Request
func (h *Client) createClient(ctx context.Context) *Request {
	var requestOpt RequestOption
	c := h.getRClient(&requestOpt)
	if h.IsBiz {
		c.OnBeforeRequest(CommonHeader(ctx))
	}
	if h.IsTrace {
		c.OnAfterResponse(ResponseTrace(ctx))
	}

	c.OnError(RequestError(ctx))
	currentParams := genDefaultParams(int(h.TimeOut), h.RetryCount, h.IdleConnTimeout, h.MaxIdleConnsPerHost)
	httpTransport := &http.Transport{MaxConnsPerHost: currentParams["defaultMaxIdleConnsPerHost"], IdleConnTimeout: time.Duration(currentParams["idleConnTimeout"])}
	httpClient := c.SetDebug(h.Debug).
		SetRetryCount(currentParams["retryCount"]).
		SetRetryMaxWaitTime(time.Duration(currentParams["timeOut"]) * time.Millisecond).
		AddRetryCondition(retryCondition()).
		SetTransport(httpTransport).
		R().SetContext(ctx)
	httpClient.SetHeader("User-Agent", h.Ua)

	return httpClient
}

// Post
//  @Description: post请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持Header,Params,Datas,Files,Auth {username,password}
//  @Return resp 请求响应
//  @Return err
func Post(ctx context.Context, url string, args ...interface{}) (resp *Response, err error) {
	req := Requests(ctx)
	resp, err = req.Post(url, args...)
	return resp, err
}

// PostJson
//  @Description: PostJson请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持Header,Params,Auth {username,password}，struct，String
//  @Return resp 请求响应
//  @Return err
func PostJson(ctx context.Context, url string, args ...interface{}) (*Response, error) {
	req := Requests(ctx)
	resp, err := req.PostJson(url, args...)
	return resp, err
}

// Get
//  @Description: get请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持Header,Params,Auth {username,password}
//  @Return resp 请求响应
//  @Return err
func Get(ctx context.Context, url string, args ...interface{}) (*Response, error) {
	req := Requests(ctx)
	resp, err := req.Get(url, args...)
	return resp, err
}

// Get
//  @Description: post请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持Header,Params,Auth {username,password}
//  @Return *Response 请求响应
//  @Return error
func (h *Client) Get(url string, args ...interface{}) (*Response, error) {
	h.httpReq = h.createClient(h.ctx)
	h.httpReq.SetHeader("Content-Type", jsonContentType)
	var params []map[string]string
	for _, arg := range args {
		switch a := arg.(type) {
		//设置请求header信息
		case Header:
			for k, v := range a {
				h.httpReq.SetHeader(k, v)
			}
		case Params:
			params = append(params, a)
		case Auth:
			// a{username,password}
			h.httpReq.SetBasicAuth(a[0], a[1])
		case Token:
			h.httpReq.SetAuthToken(string(arg.(Token)))
		default:
			h.httpReq.SetContext(h.ctx).SetBody(arg)
		}
	}

	for _, paramValue := range params {
		h.httpReq.SetQueryParams(paramValue)
	}
	h.ClientSetCookies()
	resp, err := h.httpReq.Get(url)
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

// Post
//  @Description: post请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持 Header, Params, Datas, Files, Auth {username,password}
//  @Return *Response 请求响应
//  @Return error
func (h *Client) Post(url string, args ...interface{}) (*Response, error) {
	h.httpReq = h.createClient(h.ctx)
	h.httpReq.SetHeader("Content-Type", formContentType)
	var params []map[string]string
	var datas []map[string]string // POST
	var files []map[string]string //post file

	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			for k, v := range a {
				h.httpReq.SetHeader(k, v)
			}
		case Params:
			params = append(params, a)
		case Datas: //Post form data,packaged in body.
			datas = append(datas, a)
		case Files:
			files = append(files, a)
		case Auth:
			// a{username,password}
			h.httpReq.SetBasicAuth(a[0], a[1])
		case Token:
			h.httpReq.SetAuthToken(string(arg.(Token)))
		}
	}
	h.ClientSetCookies()
	for _, paramValue := range params {
		h.httpReq.SetQueryParams(paramValue)
	}

	for _, dataValue := range datas {
		h.httpReq.SetFormData(dataValue)
	}

	for _, file := range files {
		h.httpReq.SetFiles(file)
	}

	//resp := &Response{}
	resp, err := h.httpReq.Post(url)
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

// PostJson
//  @Description: PostJson请求
//  @Param ctx 上下文
//  @Param url 请求url
//  @Param args 请求参数类型 支持 Header, Params, Auth {username,password}，struct，String
//  @Return *Response 请求响应
//  @Return error
func (h *Client) PostJson(url string, args ...interface{}) (*Response, error) {
	h.httpReq = h.createClient(h.ctx)
	h.httpReq.SetHeader("Content-Type", jsonContentType).SetContext(h.ctx)
	var params []map[string]string

	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			for k, v := range a {
				h.httpReq.SetHeader(k, v)
			}
		case Token:
			h.httpReq.SetAuthToken(string(arg.(Token)))
		case string:
			h.httpReq.SetContext(h.ctx).SetBody(arg.(string))
		case Params:
			params = append(params, a)
		case Auth:
			// a{username,password}
			h.httpReq.SetBasicAuth(a[0], a[1])
		default:
			h.httpReq.SetContext(h.ctx).SetBody(arg)
		}
	}
	h.ClientSetCookies()
	for _, paramValue := range params {
		h.httpReq.SetQueryParams(paramValue)
	}

	resp, err := h.httpReq.Post(url)
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

// getRClient
//  @Description: 获取一个新的请求client
//  @Receiver h
//  @Param requestOpt
//  @Return *RClient
func (h *Client) getRClient(requestOpt *RequestOption) *RClient {
	if requestOpt.hc != nil {
		return NewWithClient(requestOpt.hc)
	}
	return New()
}

// SetCookie
//  @Description:
//  @Receiver h
//  @Param cookie
func (h *Client) SetCookie(cookie *http.Cookie) {
	h.Cookies = append(h.Cookies, cookie)
}

// ClientSetCookies
//  @Description: 设置客户端cookie
//  @Receiver h
func (h *Client) ClientSetCookies() {
	if len(h.Cookies) > 0 {
		h.httpReq.SetCookies(h.Cookies)
	}
}

// genDefaultParams
//  @Description: 生成默认超时和重试参数
//  @Param timeOut
//  @Param retryCount
//  @Return int
//  @Return int
func genDefaultParams(timeOut, retryCount, idleConnTimeout, maxIdleConnsPerHost int) map[string]int {
	DefaultParams := map[string]int{}
	if timeOut <= 0 {
		DefaultParams["timeOut"] = defaultTimeout
	}
	if retryCount <= 0 {
		DefaultParams["retryCount"] = defaultRetryCount
	}
	if idleConnTimeout <= 0 {
		DefaultParams["idleConnTimeout"] = int(defaultIdleConnTimeout)
	}
	if maxIdleConnsPerHost <= 0 {
		DefaultParams["maxIdleConnsPerHost"] = defaultMaxIdleConnsPerHost
	}
	return DefaultParams
}
