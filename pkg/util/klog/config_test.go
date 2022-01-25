// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/8/13 2:47 下午

package klog

import (
	"testing"

	"git.bbobo.com/framework/tabby/pkg/conf"
	yaml "gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

var (
	raw = `version: v2
logging:
  # 变量值
  property:
    defaultFormat: "appId=%{appId} serviceName=%{serviceName} serviceSource=%{serviceSource} traceId=%{traceId} fileName=%{fileName} line=%{line} requestIp=%{requestIp} requestUri=%{requestUri}"
  default:
    # 和output 字段对应
    loggerType: "running"
    async: true
    level: "debug"
    #  输出源配置，可以输出多个输出源
    output:
      console:
        #format: "logging.property.defaultFormat"
        format: ""
        level: "debug"
        async: false
        # 刷新到输出源 间隔时间，单位为秒
        flushInterval: "5s"
        # 缓冲区大小
        bufferSize: 262144
`
)

type AnyFields struct {
	A string `json:"a"`
}

type myDataSource struct {
	Content string
	changed chan struct{}
}

func (m myDataSource) ReadConfig() ([]byte, error) {
	return []byte(raw), nil
}

func (m myDataSource) IsConfigChanged() <-chan struct{} {
	return m.changed
}

func (m myDataSource) Close() error {
	return nil
}

func TestConfig_Build(t *testing.T) {
	if assert.NoError(t, conf.LoadFromConfigSource(&myDataSource{}, yaml.Unmarshal)) {
		var config Config
		if assert.NoError(t, conf.UnmarshalKey("logging.default", &config)) {
			logger := config.Build()
			logger.Debug("test", Any("any", &AnyFields{A: "test"}))
		}
	}
}
