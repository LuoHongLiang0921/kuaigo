// @Description 上下文common 工具
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package klog

import (
	"context"

	"github.com/pborman/uuid"
)

type (
	ctxMarker     struct{}
	ctxLoggerType struct{}
)

var (
	ctxMarkerKey = &ctxMarker{}
	ctxLoggerKey = &ctxLoggerType{}
)

// ContextOption 上下文配置
type ContextOption struct {
	AppID         int
	ServiceName   string
	ServiceSource string
}

type Option func(c *ContextOption)

func mergeContextOption(opts ...Option) *ContextOption {
	s := &ContextOption{}
	for _, f := range opts {
		f(s)
	}
	return s
}

// WithAppID
// 	@Description
//	@Param appID
// 	@Return Option
func WithAppID(appID int) Option {
	return func(c *ContextOption) {
		c.AppID = appID
	}
}

// WithServiceName
// 	@Description 设置服务名
//	@Param serviceName 服务名
// 	@Return Option 服务名配置
func WithServiceName(serviceName string) Option {
	return func(c *ContextOption) {
		c.ServiceName = serviceName
	}
}

// WithServiceSource
// 	@Description 设置服务来源
//	@Param serviceSource 服务来源
// 	@Return Option 服务来源配置
func WithServiceSource(serviceSource string) Option {
	return func(c *ContextOption) {
		c.ServiceSource = serviceSource
	}
}

// WithCommonLog
// 	@Description: 日志上下文
//	@Param ctx 上下文
//	@Param v common 结构体
// 	@Return context.Context 带有common 的上下文 c1
func WithCommonLog(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, v)
}

// FromContext
// 	@Description: 从上下文中取common
//	@Param ctx 上下文
// 	@Return Common 上下文中common
// 	@Return bool false:上下文中存在common 日志 true:存在common 日志
func FromContext(ctx context.Context) (Common, bool) {
	if ctx == nil {
		return Common{}, false
	}
	v, ok := ctx.Value(ctxMarkerKey).(Common)
	return v, ok
}

// NewLoggerContext
// 	@Description 设置日志类型上下文
//	@Param ctx 上下文
//	@Param v 值
// 	@Return context.Context 上下文
func newLoggerContext(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, v)
}

// FromLoggerContext
// 	@Description 从日志类型上下文获取具体的日志类型
//	@Param ctx
// 	@Return string
// 	@Return bool
func FromLoggerContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	v, ok := ctx.Value(ctxLoggerKey).(string)
	return v, ok
}

// RunningLoggerContext
// 	@Description running context
//	@Param ctx 上下文
// 	@Return context.Context 运行类型的上下文
func RunningLoggerContext(ctx context.Context) context.Context {
	return newLoggerContext(ctx, LogTypeRunning)
}

// ErrorLoggerContext
// 	@Description error context
//	@Param ctx 上下文
// 	@Return context.Context 错误类型的上下文
func ErrorLoggerContext(ctx context.Context) context.Context {
	return newLoggerContext(ctx, LogTypeError)
}

// TaskLoggerContext
// 	@Description 设置任务类型context
//	@Param ctx 上下文
// 	@Return context.Context 任务类型后的上下文
func TaskLoggerContext(ctx context.Context, opts ...Option) context.Context {
	ctx = newLoggerContext(ctx, LogTypeTask)
	ctxOption := mergeContextOption(opts...)
	if com, ok := FromContext(ctx); ok {
		if com.TraceId == "" {
			com.TraceId = uuid.New()
		}
		if ctxOption.AppID != 0 {
			com.AppId = ctxOption.AppID
		}
		if ctxOption.ServiceName != "" {
			com.ServiceName = ctxOption.ServiceName
		}
		ctx = WithCommonLog(ctx, com)
		return ctx
	}
	return WithCommonLog(ctx, Common{
		AppId:         ctxOption.AppID,
		ServiceName:   ctxOption.ServiceName,
		ServiceSource: ServiceSourceTask,
	})
}

// AccessLoggerContext
// 	@Description 设置访问日志类型context
//	@Param ctx 上下文
// 	@Return context.Context
func AccessLoggerContext(ctx context.Context) context.Context {
	return newLoggerContext(ctx, LogTypeAccess)
}

func TabbyLoggerContext(ctx context.Context) context.Context {
	return newLoggerContext(ctx, LogTypeTabby)
}
