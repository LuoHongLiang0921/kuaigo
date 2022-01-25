// @Description

package klog

import (
	"context"
	"testing"
)

func Test_With_Common(t *testing.T) {
	var cfg Config
	cfg.Level = DebugLevel.String()
	cfg.Async = true

	cfg.LoggerType = LogTypeRunning
	cfg.ServiceName = "test"
	common := Common{
		AppId:         89,
		TraceId:       "900000",
		ServiceSource: "t",
		ServiceName:   "test",
		RequestIp:     "172",
		RequestUri:    "http",
		ProcessCode:   ProcessCodeRequest,
	}
	//xlog.DefaultLogger = cfg.Build()
	ctx := WithCommonLog(context.Background(), common)
	RunningLogger.WithContext(ctx).Info("test default err")
	//xlog.DefaultLogger.Info(ctx, "test default info 000")
	//xlog.DefaultLogger.Warn(ctx, "test default warn")
	//xlog.DefaultLogger.Debug(ctx, "test default debug")
	//xlog.DefaultLogger.Panic(ctx, "test default panic")
	//xlog.DefaultLogger.Fatal(ctx, "test default fatal", xlog.FieldType(xlog.LogTypeBiz), xlog.Namespace("test"))
	RunningLogger.Flush()
}

func Test_Todo(t *testing.T) {
	//t1()
	var cfg Config
	cfg.Level = DebugLevel.String()
	cfg.Store = []string{"file", "console"}

	cfg.Async = true
	cfg.CallerSkip = 1
	cfg.AddCaller = true
	cfg.LoggerType = LogTypeAccess
	cfg.ServiceName = "test_todo"
	//
	//common := xlog.Common{
	//	AppId:         89,
	//	TraceId:       "900000",
	//	ServiceSource: "t",
	//	ServiceName:   "test",
	//	RequestIp:     "172",
	//	RequestUri:    "http",
	//	ProcessCode:   xlog.ProcessCodeRequest,
	//}
	RunningLogger = cfg.WithConfigVersion("v2").Build()
	ctx := context.TODO()
	RunningLogger.WithContext(ctx).Error("test default err")
	RunningLogger.WithContext(ctx).Info("test default info")
	RunningLogger.WithContext(ctx).Warn("test default warn")
	RunningLogger.WithContext(ctx).Debug("test default debug")
	RunningLogger.WithContext(ctx).Panic("test default panic")
	RunningLogger.WithContext(ctx).Fatal("test default fatal")
	RunningLogger.Flush()
}
