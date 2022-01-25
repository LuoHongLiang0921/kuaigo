// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/6/29 10:35 上午

package rotate

import (
	"io"
	"time"
)

type RotateOption func(l *Logger)

// WithFileName 设置文件名
// 	@Description
//	@Param name
// 	@Return RotateOption
func WithFileName(name string) RotateOption {
	return func(l *Logger) {
		l.Filename = name
	}
}

// WithMaxSize 设置文件最大大小，单位MB
// 	@Description
//	@Param maxSize
// 	@Return RotateOption
func WithMaxSize(maxSize int) RotateOption {
	return func(l *Logger) {
		l.MaxSize = maxSize
	}
}

// WithMaxAge 设置文件年龄 d
// 	@Description
//	@Param maxAge
// 	@Return RotateOption
func WithMaxAge(maxAge int) RotateOption {
	return func(l *Logger) {
		l.MaxAge = maxAge
	}
}

func WithMaxBackups(maxBackups int) RotateOption {
	return func(l *Logger) {
		l.MaxBackups = maxBackups
	}
}

func WithInterval(interval time.Duration) RotateOption {
	return func(l *Logger) {
		l.Interval = interval
	}
}

// NewRotateWriter
// 	@Description 文件切分配置
//	@Param config
// 	@Return io.Writer
func NewRotateWriter(opts ...RotateOption) io.Writer {
	rotateLog := NewLogger()
	rotateLog.LocalTime = true
	rotateLog.Compress = false
	for _, f := range opts {
		f(rotateLog)
	}
	//rotateLog.Filename = config.Filename()
	//rotateLog.MaxSize = config.MaxSize // MB
	//rotateLog.MaxAge = config.MaxAge   // days
	//rotateLog.MaxBackups = config.MaxBackup
	//rotateLog.Interval = config.Interval

	return rotateLog
}
