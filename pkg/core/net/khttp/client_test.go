package khttp

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/http"
	"testing"
)

type requestBody struct {
	BusinessId int    `json:"businessId"`
	Content    string `json:"content"`
	ContentId  int    `json:"contentId"`
	ReplyId    int    `json:"replyId"`
	RootId     int    `json:"rootId"`
	ToUserId   int    `json:"toUserId"`
	UserId     int    `json:"userId"`
	Resources  string `json:"resources"`
	Ext        string `json:"ext"`
}

var res = &requestBody{
	BusinessId: 56941891,
	Content:    "est reprehenderit tempor adipisicing",
	ContentId:  64062154,
	ReplyId:    7475653,
	RootId:     15640399,
	ToUserId:   76784414,
	UserId:     30170306,
	Resources:  "minim ipsum eu adipisicing",
	Ext:        "non voluptate",
}

func TestClient_Post(t *testing.T) {
	//https://httpbin.org/post
	//http://192.168.132.70:30633//service/interaction/comment/getCommentCount
	oriUrl := "http://192.168.132.70:10633//service/interaction/comment/createComment"
	resp, err := Post(context.Background(), oriUrl, Header{"ai": "101"}, Datas{"AC": "123", "AA": "456"})
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
		return
	}
	type jsonResp struct {
		Code int
		Msg  string
	}
	var json jsonResp
	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	resp.Json(&json)
	fmt.Println("json 结构")
	fmt.Printf("%#v\n", json)
	klog.FlushAll()
}

// TestClient_Request_Get
//  @Description: 高级方式get
//  @Param t
func TestClient_Request_Get(t *testing.T) {
	req := Requests(context.Background())
	req.Debug = true
	req.IsBiz = false
	oriUrl := "https://api.apiopen.top/getJoke"
	req.SetCookie(&http.Cookie{
		Name:  "go-resty",
		Value: "This is cookie value",
	})
	resp, err := req.Get(oriUrl, Header{"ai": "101"}, Params{"page": "1", "count": "2", "type": "video"})
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}
	type jsonResp struct {
		Code   int
		Msg    string
		Result interface{}
	}
	var json jsonResp
	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	resp.Json(&json)
	fmt.Println("json 结构")
	fmt.Printf("%#v\n", json)

}

// TestClient_Trace
//  @Description: 输出响应trace信息
//  @Param t
func TestClient_Trace(t *testing.T) {
	req := Requests(context.Background())

	req.Debug = true
	req.IsBiz = false
	req.IsTrace = true
	oriUrl := "https://api.apiopen.top/getJoke"
	req.SetCookie(&http.Cookie{
		Name:  "go-resty",
		Value: "This is cookie value",
	})

	resp, err := req.Get(oriUrl, Header{"ai": "101", "test": "123"}, Params{"page": "1", "count": "2", "type": "video"})
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}

	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

func TestClient_Request_GetBiz(t *testing.T) {
	ctx := klog.WithCommonLog(context.Background(), klog.Common{})
	req := Requests(ctx)
	req.Debug = true
	req.IsBiz = true
	oriUrl := "https://api.apiopen.top/getJoke"
	resp, err := req.Get(oriUrl, Header{"ai": "102", "ti": "test"}, Params{"page": "1", "count": "2", "type": "video"})
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}
	type jsonResp struct {
		Code   int
		Msg    string
		Result interface{}
	}
	var json jsonResp
	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	resp.Json(&json)
	fmt.Println("json 结构")
	fmt.Printf("%#v\n", json)

}

func TestClient_Get(t *testing.T) {
	oriUrl := "https://api.apiopen.top/getJoke"
	resp, err := Get(context.Background(), oriUrl, Header{"ai": "101"}, Params{"page": "1", "count": "2", "type": "video"})
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}
	type jsonResp struct {
		Code int
		Msg  string
	}
	var json jsonResp
	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	resp.Json(&json)
	fmt.Println("json 结构")
	fmt.Printf("%#v\n", json)

}

// TestClient_PostJsonWithToken
//  @Description: 携带authtoken请求
//  @Param t
func TestClient_PostJsonWithToken(t *testing.T) {
	oriUrl := "https://httpbin.org/post"
	resp, err := PostJson(context.Background(), oriUrl, Header{"ai": "101"}, res, Token("1233344"))
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}

	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	//resp.Json(&json)
	//fmt.Println("json 结构")
	//fmt.Printf("%#v\n", json)
}

// TestClient_PostJson
//  @Description: post方式提交json
//  @Param t
func TestClient_PostJson(t *testing.T) {
	//https://httpbin.org/post
	//http://192.168.132.70:30633//service/interaction/comment/getCommentCount
	oriUrl := "http://192.168.132.70:30633//service/interaction/comment/createComment"

	resp, err := PostJson(context.Background(), oriUrl, Header{"ai": "101"}, res)
	if err != nil {
		klog.Error("post2 ", klog.FieldErr(err))
	}
	type jsonResp struct {
		Code int
		Msg  string
	}
	var json jsonResp
	fmt.Println("=======================")
	fmt.Println("响应string", resp.String())
	fmt.Println("响应size", resp.Size())
	fmt.Println("响应StatusCode", resp.StatusCode())
	//var json map[string]interface{}
	resp.Json(&json)
	fmt.Println("json 结构")
	fmt.Printf("%#v\n", json)

}
