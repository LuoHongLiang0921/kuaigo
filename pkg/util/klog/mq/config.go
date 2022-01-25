// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/6/29 5:10 下午

package mq

import (
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Level 日志等级
	Level string
	//
	Source string
}

func (c Config) Build() zapcore.Core {
	panic("implement me")
}
