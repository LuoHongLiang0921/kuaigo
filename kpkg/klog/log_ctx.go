package klog

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/defers"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kcolor"
	"regexp"

	"github.com/pborman/uuid"

	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	FatalLevel = zap.FatalLevel
)

// Func ...
type (
	Func  func(context.Context, string, ...zap.Field)
	Field = zap.Field
	Level = zapcore.Level
)

type Logger struct {
	desugar    *zap.Logger
	lv         *zap.AtomicLevel
	config     Config
	sugar      *zap.SugaredLogger
	loggerType string
}

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
	Skip = zap.Skip()
	// ByteString ...
	ByteString = zap.ByteString
)
var gormSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	gormSourceDir = regexp.MustCompile(`log_ctx\.go`).ReplaceAllString(file, "")
}

func FileWithLineNum() (string, int) {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file, line
			}
			idx = strings.LastIndexByte(file[:idx], '/')
			if idx == -1 {
				return file, line
			}
			return file[idx+1:], line
		}
	}
	return "", 0
}

func newLogger(config *Config) *Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}

	var ws zapcore.WriteSyncer
	if config.Debug {
		ws = os.Stdout
	} else {
		ws = zapcore.AddSync(newRotate(config))
	}

	if config.Async {
		var close CloseFunc
		ws, close = Buffer(ws, defaultBufferSize, defaultFlushInterval)
		defers.Register(close)
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}
	encoderConfig := *config.EncoderConfig
	redisZapCore := config.RedisCore
	// 加入redis log 输出流
	if redisZapCore == nil && config.redisClient != nil && config.Save {
		//rws := zapcore.AddSync(redis.NewRedisLog(config.LogKey, config.redisClient))
		//encoder := func() zapcore.Encoder {
		//	return zapcore.NewJSONEncoder(encoderConfig)
		//}
		//redisZapCore = zapcore.NewCore(encoder(), rws, lv)
	}
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.Debug {
					return zapcore.NewConsoleEncoder(*ConsoleZapConfig())
				}
				return zapcore.NewJSONEncoder(encoderConfig)
			}(),
			ws,
			lv,
		)
	}
	var cores []zapcore.Core
	if core != nil {
		cores = append(cores, core)
	}
	if redisZapCore != nil {
		cores = append(cores, redisZapCore)
	}
	zapLogger := zap.New(
		zapcore.NewTee(cores...),
		zapOptions...,
	)
	return &Logger{
		desugar:    zapLogger,
		lv:         &lv,
		config:     *config,
		sugar:      zapLogger.Sugar(),
		loggerType: config.LoggerType,
	}
}

func (logger *Logger) getContext() context.Context {
	return context.TODO()
}

// AutoLevel ...
func (logger *Logger) AutoLevel(confKey string) {
	conf.OnChange(func(config *conf.Configuration) {
		lvText := strings.ToLower(config.GetString(confKey))
		if lvText != "" {
			logger.Info(logger.getContext(), "update level", String("level", lvText), String("name", logger.config.Name))
			logger.lv.UnmarshalText([]byte(lvText))
		}
	})
}

// SetLevel ...
func (logger *Logger) SetLevel(lv Level) {
	logger.lv.SetLevel(lv)
}

func (logger *Logger) SetServiceName(s string) {
	logger.config.ServiceName = s
}

// Flush ...
func (logger *Logger) Flush() error {
	return logger.desugar.Sync()
}

// ConsoleZapConfig ...
func ConsoleZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "logLevel",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		ConsoleSeparator:" ",
	}
}

func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       zapcore.OmitKey,
		NameKey:        zapcore.OmitKey,
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// DebugEncodeLevel ...
func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = kcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = kcolor.Blue
	case zapcore.InfoLevel:
		colorize = kcolor.Green
	case zapcore.WarnLevel:
		colorize = kcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = kcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 毫秒
	enc.AppendInt64(t.UnixNano() / 1e6)
}

// IsDebugMode ...
func (logger *Logger) IsDebugMode() bool {
	return logger.config.Debug
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// StdLog ...
func (logger *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(logger.desugar)
}

func (logger *Logger) makeFields(ctx context.Context, level Level) []Field {
	com, ok := FromContext(ctx)
	var fields []Field
	file, line := FileWithLineNum()
	if !ok {
		com = Common{}
	}
	traceId := com.TraceId
	if traceId == "" {
		traceId = uuid.New()
	}
	fields = append(fields, Int("appId", com.AppId))
	fields = append(fields, String("traceId", traceId))
	fields = append(fields, String("serviceSource", com.ServiceSource))
	fields = append(fields, String("serviceName", logger.config.ServiceName))
	switch logger.loggerType {
	case LogTypeError:
		fallthrough
	case LogTypeRunning:
		fields = append(fields, String("fileName", file))
		fields = append(fields, String("logLevel", level.String()))
		fields = append(fields, Int("line", line))
		fields = append(fields, String("requestIp", com.RequestIp))
		fields = append(fields, String("requestUri", com.RequestUri))
	case LogTypeAccess:
		fields = append(fields, String("requestIp", com.RequestIp))
		fields = append(fields, String("requestUri", com.RequestUri))
		fields = append(fields, Int("processCode", com.ProcessCode))
		fields = append(fields, Int("costTime", int(com.CostTime)))
		fields = append(fields, Int("code", com.Code))
		fields = append(fields, Any("p", com.P))
		fields = append(fields, String("uid", com.UID))
		fields = append(fields, String("msg", ""))
		fields = append(fields, String("logLevel", ""))
	case LogTypeTask:
	}
	//return fields
	//traceId := uuid.New()
	//appID := 0
	//serverSource := ""
	//fields = append(fields, Int("appId", appID))
	//fields = append(fields, String("traceId", traceId))
	//fields = append(fields, String("serviceSource", serverSource))
	//fields = append(fields, String("serviceName", logger.config.ServiceName))
	return fields

}

// Debug ...
func (logger *Logger) Debug(ctx context.Context, msg string, fields ...Field) {

	//if logger.IsDebugMode() {
	//	msg = normalizeMessage(msg)
	//}
	fss := logger.makeFields(ctx, DebugLevel)
	fss = append(fss, fields...)
	logger.desugar.Debug(msg, fss...)
}

// Debugf ...
func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.Debug(context.TODO(), sprintf(template, args...))
}

// Info ...
func (logger *Logger) Info(ctx context.Context, msg string, fields ...Field) {
	fss := logger.makeFields(ctx, InfoLevel)
	fss = append(fss, fields...)
	logger.desugar.Info(msg, fss...)
}

// Infof ...
func (logger *Logger) Infof(template string, args ...interface{}) {
	//logger.sugar.Infof(sprintf(template, args...))
	logger.Info(context.TODO(), sprintf(template, args...))
}

// Warn ...
func (logger *Logger) Warn(ctx context.Context, msg string, fields ...Field) {
	//if logger.IsDebugMode() {
	//	msg = normalizeMessage(msg)
	//}
	fss := logger.makeFields(ctx, WarnLevel)
	fss = append(fss, fields...)
	logger.desugar.Warn(msg, fss...)
}

// Warnf ...
func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.Warn(context.TODO(), sprintf(template, args...))
}

// Error ...
// todo: 可以使用zap.Object(),性能会高一些
func (logger *Logger) Error(ctx context.Context, msg string, fields ...Field) {
	//if logger.IsDebugMode() {
	//	msg = normalizeMessage(msg)
	//}
	fss := logger.makeFields(ctx, ErrorLevel)
	fss = append(fss, fields...)
	logger.desugar.Error(msg, fss...)
}

//
// Errorf ...
func (logger *Logger) Errorf(template string, args ...interface{}) {
	//logger.sugar.Errorf(sprintf(template, args...))
	logger.Error(context.TODO(), sprintf(template, args...))
}

// Panic ...
func (logger *Logger) Panic(ctx context.Context, msg string, fields ...Field) {
	fss := logger.makeFields(ctx, PanicLevel)
	fss = append(fss, fields...)
	logger.desugar.Panic(msg, fss...)
}

func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.Panic(context.TODO(), sprintf(template, args...))
}

// DPanic ...
func (logger *Logger) DPanic(ctx context.Context, msg string, fields ...Field) {
	fss := logger.makeFields(ctx, zap.DPanicLevel)
	fss = append(fss, fields...)
	logger.desugar.DPanic(msg, fss...)
}

func (logger *Logger) DPanicf(template string, args ...interface{}) {
	logger.DPanic(context.TODO(), sprintf(template, args...))
}

// Fatal ...
func (logger *Logger) Fatal(ctx context.Context, msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
		return
	}
	fss := logger.makeFields(ctx, FatalLevel)
	fss = append(fss, fields...)
	logger.desugar.Fatal(msg, fss...)
}

func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.Fatal(context.TODO(), sprintf(template, args...))
}

func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	// 控制台输出
	fmt.Printf("%s: \n    %s: %s\n", kcolor.Red("panic"), kcolor.Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", kcolor.Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", kcolor.Red(key), fmt.Sprintf("%+v", val))
	}

}

// With ...
// todo: zap.Namespace
func (logger *Logger) With(fields ...Field) *Logger {
	desugarLogger := logger.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}
