// @Description
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package klog

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
)

type Func func(string, ...zap.Field)

var (
	KuaigoLogger   = getDefaultConfig("").WithConfigVersion("v2").WithLoggerType(LogTypeTabby).RegisterOutput().Build()
	RunningLogger = getDefaultConfig("").WithConfigVersion("v2").WithLoggerType(LogTypeRunning).RegisterOutput().Build()
	ErrorLogger   = getDefaultConfig("").WithConfigVersion("v2").WithLoggerType(LogTypeError).RegisterOutput().Build()
	AccessLogger  = getDefaultConfig("").WithConfigVersion("v2").WithLoggerType(LogTypeAccess).RegisterOutput().Build()
	TaskLogger    = getDefaultConfig("").WithConfigVersion("v2").WithLoggerType(LogTypeTask).RegisterOutput().Build()

	defaultLogger = RunningLogger
)

// Auto
// 	@Description RunningLogger err 不为nil时，为Error 级别，否则是Info 级别
//	@Param err 错误信息
// 	@return Func 日志函数
func Auto(err error) Func {
	if err != nil {
		return defaultLogger.With(zap.Any("err", err.Error())).Error
	}
	return defaultLogger.Info
}

// Info
// 	@Description  运行时 Info 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func Error(msg string, fields ...zap.Field) {
	ErrorLogger.Error(msg, fields...)
}

// Panic
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg字段内容
//	@Param fields 自定义字段内容
func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}

// DPanic
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg字段内容
//	@Param fields 自定义字段内容
func DPanic(msg string, fields ...zap.Field) {
	defaultLogger.DPanic(msg, fields...)
}

// Fatal
// 	@Description 运行时 Debug 级别日志输出
//	@Param ctx 上下文
//	@Param msg msg字段内容
//	@Param fields 自定义字段内容
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Debugf 运行时Debug级别日志
// 	@Description RunningLogger 运行时 Debug 级别日志格式化输出
//	@Param template  msg字段模板内容
//	@Param args 格式化字符对应的值
func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

// Infof
// 	@Description 运行时 Info 级别日志格式化输出
//	@Param template msg字段格式字符
//	@Param args 格式化字符对应的值
func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

// Warnf
// 	@Description 运行时 warn 级别日志格式化输出
//	@Param template msg字段格式字符 msg字段内容
//	@Param args 格式化字符对应的值 格式化字符对应的值
func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

// Errorf
// 	@Description 错误日志类型 error 级别日志格式化输出
//	@Param template msg字段格式字符
//	@Param args 格式化字符对应的值
func Errorf(template string, args ...interface{}) {
	ErrorLogger.Errorf(template, args...)
}

// Panicf
// 	@Description 运行时 panic 级别日志格式化输出
//	@Param template msg字段格式字符
//	@Param args 格式化字符对应的值
func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

// DPanicf
// 	@Description 运行时 debug panic 级别日志输出
//	@Param template msg字段格式字符
//	@Param args 格式化字符对应的值
func DPanicf(template string, args ...interface{}) {
	defaultLogger.DPanicf(template, args...)
}

// Fatalf
// 	@Description 运行时 fatal 级别日志格式化输出
//	@Param template msg字段格式字符 msg字段内容
//	@Param args 格式化字符对应的值
func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

// Log
// 	@Description 日志函数
// 	@receiver fn
//	@Param ctx 上下文
//	@Param msg msg字段内容
//	@Param fields 自定义字段内容
func (fn Func) Log(msg string, fields ...zap.Field) {
	fn(msg, fields...)
}

// With
// 	@Description 运行时 直接使用zap 日志库
//	@Param fields 自定义字段内容
// 	@return *Logger
func With(fields ...zap.Field) *Logger {
	return defaultLogger.With(fields...)
}

// WithContext
// 	@Description 设置一个
//	@Param ctx 上下文
//	@Param keys context 中key 列表
// 	@Return *Logger 带有上下文后的日志实例
func WithContext(ctx context.Context, keys ...interface{}) *Logger {
	return defaultLogger.WithContext(ctx, keys...)
}

// FlushAll
// @Description 刷新所有日志
func FlushAll() {
	// todo: so ugly，统一处理
	_ = RunningLogger.Flush()
	_ = KuaigoLogger.Flush()
	_ = ErrorLogger.Flush()
	_ = AccessLogger.Flush()
	_ = TaskLogger.Flush()
}

// InitLogger
// 	@Description 初始化日志
//	@Param configVersion 配置版本
//	@Param serviceName 服务或应用名字
func InitLogger(configVersion, serviceName string) {
	runningKey := "running"
	tabbyKey := "tabbyLogger"
	accessKey := "access"
	errKey := "error"
	taskKey := "task"
	if configVersion > "" {
		runningKey = "logging.running"
		tabbyKey = "logging.default"
		accessKey = "logging.access"
		errKey = "logging.error"
		taskKey = "logging.task"
	}
	if conf.Get(runningKey) != nil {
		RunningLogger = RawConfig(runningKey).WithConfigVersion(configVersion).Build().clone()
	}
	RunningLogger.SetServiceName(serviceName)

	if conf.Get(tabbyKey) != nil {
		KuaigoLogger = RawConfig(tabbyKey).WithConfigVersion(configVersion).Build().clone()
	}
	KuaigoLogger.SetServiceName(serviceName)

	if conf.Get(accessKey) != nil {
		AccessLogger = RawConfig(accessKey).WithConfigVersion(configVersion).Build().clone()
	}
	AccessLogger.SetServiceName(serviceName)

	if conf.Get(errKey) != nil {
		ErrorLogger = RawConfig(errKey).WithConfigVersion(configVersion).Build().clone()
	}
	ErrorLogger.SetServiceName(serviceName)

	if conf.Get(taskKey) != nil {
		TaskLogger = RawConfig(taskKey).WithConfigVersion(configVersion).Build().clone()
	}
	TaskLogger.SetServiceName(serviceName)

}
