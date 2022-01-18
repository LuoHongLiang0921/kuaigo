package klog

import (
	"context"

	"go.uber.org/zap"
)

const (
	RedisURL = "dev-logproxy.yixiahd.com:6379"
)

// frame logger
var KuaigoLogger = Config{
	Debug:      true,
	Save:       true,
	LoggerType: LogTypeRunning,
	RedisURL:   RedisURL,
	LogKey:     "running",
}.Build()

//RunningLogger 运行时日志
var RunningLogger = Config{
	Debug:      true,
	Save:       true,
	Async:      true,
	LoggerType: LogTypeRunning,
	RedisURL:   RedisURL,
	LogKey:     "running",
}.Build()

//ErrorLogger 错误日志
var ErrorLogger = Config{
	Debug:      true,
	Save:       true,
	Async:      true,
	LoggerType: LogTypeError,
	RedisURL:   RedisURL,
	LogKey:     "error",
}.Build()

//AccessLogger 访问日志
var AccessLogger = Config{
	Debug:      true,
	Save:       true,
	Async:      true,
	LoggerType: LogTypeAccess,
	RedisURL:   RedisURL,
	LogKey:     "access",
}.Build()

//TaskLogger 任务日志
var TaskLogger = Config{
	Debug:      true,
	Save:       true,
	LoggerType: LogTypeTask,
	RedisURL:   RedisURL,
	LogKey:     "task",
}.Build()

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return RunningLogger.With(zap.Any("err", err.Error())).Error
	}

	return RunningLogger.Info
}

// Info ...
func Info(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Info(ctx, msg, fields...)
}

// Debug ...
func Debug(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Debug(ctx, msg, fields...)
}

// Warn ...
func Warn(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Warn(ctx, msg, fields...)
}

// Error ...
func Error(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Error(ctx, msg, fields...)
	//DefaultLogger.Error(msg, fields...)
}

// Panic ...
func Panic(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Panic(ctx, msg, fields...)
	//DefaultLogger.Panic(msg, fields...)
}

// DPanic ...
func DPanic(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.DPanic(ctx, msg, fields...)
	//DefaultLogger.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(ctx context.Context, msg string, fields ...Field) {
	RunningLogger.Fatal(ctx, msg, fields...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	RunningLogger.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	RunningLogger.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	RunningLogger.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	RunningLogger.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	RunningLogger.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	RunningLogger.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	RunningLogger.Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(ctx context.Context, msg string, fields ...Field) {
	fn(ctx, msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return RunningLogger.With(fields...)
}
