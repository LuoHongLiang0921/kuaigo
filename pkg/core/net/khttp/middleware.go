// @Description  常用中间件

package khttp

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"strconv"

	"github.com/go-resty/resty/v2"

	"github.com/google/uuid"
)

// Metric 指标监控
// 	@Description: promethus
//	@Param ctx 上下文
//	@Param name 服务名字
// 	@Return AfterResponseHandler
func Metric(ctx context.Context, name string) ResponseMiddleware {
	return func(c *RClient, req *resty.Response) error {
		return nil
	}
}

// CommonHeader
// 	@Description: 根据上下文，请求头中加入ai，ti
//	@Param ctx 上下文
// 	@Return BeforeRequestHandler
func CommonHeader(ctx context.Context) RequestMiddleware {
	return func(c *RClient, req *Request) error {
		if com, ok := klog.FromContext(ctx); ok {
			headers := req.Header
			if headers.Get(constant.HeaderFieldAi) == "" {
				req.SetHeader(constant.HeaderFieldAi, strconv.Itoa(com.AppId))
			}
			if headers.Get(constant.HeaderFieldTi) == "" {
				traceId := com.TraceId
				if traceId == "" {
					traceId = uuid.New().String()
				}
				req.SetHeader(constant.HeaderFieldTi, traceId)
			}
		}
		return nil
	}
}

// ResponseTrace
// 	@Description: 根据上下文，请求头中加入ai，ti
//	@Param ctx 上下文
// 	@Return BeforeRequestHandler
func ResponseTrace(ctx context.Context) ResponseMiddleware {
	return func(c *RClient, resp *resty.Response) error {

		logTxt := fmt.Sprintf("httpLog RequestUrl:%s StatusCode:%v costTime:%s QueryParam:%+v RequestBody:%+v", resp.Request.URL, resp.StatusCode(), resp.Time(), resp.Request.QueryParam, resp.Request.Body)
		klog.RunningLogger.WithContext(ctx).Info(logTxt)
		return nil
	}
}

func RequestError(ctx context.Context) resty.ErrorHook {
	return func(req *Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logTxt := fmt.Sprintf("httpLog RequestUrl:%s StatusCode:%v Error:%s,Response:%s", req.URL, v.Response.StatusCode(), v.Error(), v.Response.String())

			klog.ErrorLogger.WithContext(ctx).Error(logTxt)
		}
		// Log the error, increment a metric, etc...
	}
}
