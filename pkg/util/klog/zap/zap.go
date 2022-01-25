// @Description zap
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/4/22 2:11 下午

package zap

import (
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
	Skip = zap.Skip()
	// ByteString ...
	ByteString               = zap.ByteString
	Error                    = zap.Error
	NewStdLog                = zap.NewStdLog
	AddSync                  = zapcore.AddSync
	NewAtomicLevelAt         = zap.NewAtomicLevelAt
	NewCore                  = zapcore.NewCore
	NewJSONEncoder           = zapcore.NewJSONEncoder
	NewConsoleEncoder        = zapcore.NewConsoleEncoder
	DefaultLineEnding        = zapcore.DefaultLineEnding
	ShortCallerEncoder       = zapcore.ShortCallerEncoder
	FullCallerEncoder        = zapcore.FullCallerEncoder
	CapitalColorLevelEncoder = zapcore.CapitalColorLevelEncoder
	LowercaseLevelEncoder    = zapcore.LowercaseLevelEncoder
	SecondsDurationEncoder   = zapcore.SecondsDurationEncoder
	StringDurationEncoder    = zapcore.StringDurationEncoder
	MillisDurationEncoder    = zapcore.MillisDurationEncoder
	EpochMillisTimeEncoder   = zapcore.EpochMillisTimeEncoder
	NanosDurationEncoder     = zapcore.NanosDurationEncoder
	FullNameEncoder          = zapcore.FullNameEncoder
	OmitKey                  = zapcore.OmitKey
	New                      = zap.New
	NewTee                   = zapcore.NewTee
	NewDevelopmentConfig     = zap.NewDevelopmentConfig
)

// struct
type (
	Logger        = zap.Logger
	AtomicLevel   = zap.AtomicLevel
	SugaredLogger = zap.SugaredLogger

	Field                 = zap.Field
	Level                 = zapcore.Level
	Core                  = zapcore.Core
	WriteSyncer           = zapcore.WriteSyncer
	LevelEnablerFunc      = zap.LevelEnablerFunc
	Encoder               = zapcore.Encoder
	EncoderConfig         = zapcore.EncoderConfig
	PrimitiveArrayEncoder = zapcore.PrimitiveArrayEncoder
	Option                = zap.Option
)
