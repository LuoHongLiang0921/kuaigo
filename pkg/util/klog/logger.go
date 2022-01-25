package klog

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"log"

)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-Level logs.
	ErrorLevel = zap.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel  = zap.FatalLevel
	DPanicLevel = zap.DPanicLevel
)

var (
	// String ...
	String = zap.String
	// Any ...
	Any = zap.Any
	// Int64 ...
	Int64 = zap.Int64
	// Int ...
	Int = zap.Int
	// Int32 ...
	Int32 = zap.Int32
	// Uint ...
	Uint = zap.Uint
	// Duration ...
	Duration = zap.Duration
	// Durationp ...
	Durationp = zap.Durationp
	// Object ...
	Object = zap.Object
	// Namespace ...
	Namespace = zap.Namespace
	// Reflect ...
	Reflect = zap.Reflect
	// Skip ...
	Skip = zap.Skip
	// ByteString ...
	ByteString = zap.ByteString
)

type (
	Field = zap.Field
)

// Logger 日志
type Logger struct {
	desugar    *zap.Logger
	config     *Config
	sugar      *zap.SugaredLogger
	loggerType string
	ctx        context.Context
	// 是否自动构造common 字段，默认为true
	isPrintCommon bool
}

func (lg *Logger) getContext() context.Context {
	if lg.ctx == nil {
		return context.Background()
	}
	return lg.ctx
}

// WithContext
// 	@Description 设置上下文,根据上下为中的key
// 	@Receiver lg Logger
//	@Param ctx 上下文
//	@Param keys 上下文中的key 列表
// 	@Return *Logger 带有上下文的日志
func (lg *Logger) WithContext(ctx context.Context, keys ...interface{}) *Logger {
	if ctx == nil {
		return lg
	}
	logger := (*Logger)(nil)
	if loggerType, ok := FromLoggerContext(ctx); ok {
		switch loggerType {
		case LogTypeRunning:
			logger = RunningLogger
		case LogTypeAccess:
			logger = AccessLogger
		case LogTypeError:
			logger = ErrorLogger
		case LogTypeTask:
			logger = TaskLogger
		case LogTypeTabby:
			logger = KuaigoLogger
		default:
			logger = lg
		}
	} else {
		logger = lg
	}
	newLg := logger.clone()
	newLg.ctx = ctx
	return newLg
}

func (lg *Logger) clone() *Logger {
	copy := *lg
	return &copy
}

// SetServiceName
// 	@Description 设置服务名
// 	@receiver lg Logger
//	@Param s 服务名
func (lg *Logger) SetServiceName(s string) {
	lg.config.ServiceName = s
}

// Flush
// 	@Description 将buffer 中的数据写入到存储中
// 	@receiver logger Logger
// 	@return error 错误
func (lg *Logger) Flush() error {
	return lg.desugar.Sync()
}

// StdLog
// 	@Description 标准日志
// 	@receiver logger Logger
// 	@return *log.Logger zap标准日志
func (lg *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(lg.desugar)
}

// Debug
// 	@Description 日志 Debug 级别日志输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Debug(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.DebugLevel)
	fss = append(fss, fields...)
	lg.desugar.Debug(msg, fss...)
}

// Debugf
// 	@Description Debug 级别格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Debugf(template string, args ...interface{}) {
	lg.Debug(sprintf(template, args...))
}

// Info
// 	@Description 日志 info 级别日志输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Info(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.InfoLevel)
	fss = append(fss, fields...)
	lg.desugar.Info(msg, fss...)
}

// Infof
// 	@Description info 级别日志格式化输出
// 	@Receiver lg  Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Infof(template string, args ...interface{}) {
	lg.Info(sprintf(template, args...))
}

// Warn
// 	@Description warn 级别日志格式化输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Warn(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.WarnLevel)
	fss = append(fss, fields...)
	lg.desugar.Warn(msg, fss...)
}

// Warnf
// 	@Description warn 级别日志格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Warnf(template string, args ...interface{}) {
	lg.Warn(sprintf(template, args...))
}

// Error
// 	@Description error 级别日志格输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Error(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.ErrorLevel)
	fss = append(fss, fields...)
	lg.desugar.Error(msg, fss...)
}

// Errorf
// 	@Description error 级别日志格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Errorf(template string, args ...interface{}) {
	lg.Error(sprintf(template, args...))
}

// Panic
// 	@Description panic  级别日志输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Panic(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.PanicLevel)
	fss = append(fss, fields...)
	lg.desugar.Panic(msg, fss...)
}

// Panicf
//  @Deprecated  请使用 PanicMsgf
// 	@Description panic  级别日志格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Panicf(template string, args ...interface{}) {
	lg.Panic(sprintf(template, args...))
}

// DPanic
// 	@Description DPanic 级别日志输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) DPanic(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.DPanicLevel)
	fss = append(fss, fields...)
	lg.desugar.DPanic(msg, fss...)
}

// DPanicf
// 	@Description DPanic 级别日志格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) DPanicf(template string, args ...interface{}) {
	lg.DPanic(sprintf(template, args...))
}

// Fatal
// 	@Description fatal 级别日志输出
// 	@Receiver lg Logger
//	@Param msg msg 字段内容
//	@Param fields 自定义字段内容
func (lg *Logger) Fatal(msg string, fields ...zap.Field) {
	fss := lg.makeFields(zap.FatalLevel)
	fss = append(fss, fields...)
	lg.desugar.Fatal(msg, fss...)
}

func (lg *Logger) SetLevel(l zap.Level) {

}

// Fatalf
//  @Deprecated 请使用 FatalMsgf
// 	@Description fatal 级别日志格式化输出
// 	@Receiver lg Logger
//	@Param template msg字段模板内容
//	@Param args 格式化字符对应的值
func (lg *Logger) Fatalf(template string, args ...interface{}) {
	lg.Fatal(sprintf(template, args...))
}

// With
// 	@Description 直接使用日志
// 	@Receiver lg Logger
//	@Param fields 字段
// 	@return *Logger 字段构成的日志
// todo: zap.Namespace
func (lg *Logger) With(fields ...zap.Field) *Logger {
	desugarLogger := lg.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		//lv:      lg.lv,
		sugar:  desugarLogger.Sugar(),
		config: lg.config,
	}
}
