package console

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/defers"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/buffer"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/property"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"os"
	"strings"
	"time"

)

const (
	OutputConsole = "console"
)

type Config struct {

	// Level 日志初始等级
	Level string
	// 是否添加调用者信息
	AddCaller  bool
	CallerSkip int

	// 格式化样式串
	Format Format

	//Async 异步
	Async bool
	// 刷新间隔
	FlushInterval string
	BufferSize    int

	parentKey string           `mapstructure:"-"`
	lv        *zap.AtomicLevel `mapstructure:"-"`
}

// SetDefaultConfig 设置console 输出源默认配置
// 	@Description
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetDefaultConfig() *Config {
	if c.FlushInterval == "" {
		c.FlushInterval = "5s"
	}
	if c.BufferSize == 0 {
		c.BufferSize = 256 * 1024
	}
	return c
}

// GetFormatValue
// 	@Description 获取日志输出格式
// 	@Receiver c
// 	@Return string
func (c *Config) GetFormatValue() string {
	if c.Format != "" {
		return conf.GetString(string(c.Format))
	}
	return string(c.Format)
}

// SetParent
// 	@Description 设置父级key
// 	@Receiver c
//	@Param k
// 	@Return *Config
func (c *Config) SetParent(k string) *Config {
	c.parentKey = k
	return c
}

// SetAutoLevel
// 	@Description 设置auto level
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetAutoLevel() *Config {
	var lv zap.Level

	err := lv.Set(c.Level)
	if err != nil {
		panic(err)
	}
	alv := zap.NewAtomicLevelAt(lv)
	err = alv.UnmarshalText([]byte(c.Level))
	if err != nil {
		panic(err)
	}
	c.lv = &alv
	if c.parentKey == "" {
		return c
	}
	conf.OnChange(func(config *conf.Configuration) {
		lvText := strings.ToLower(config.GetString(c.getLevelVariable()))
		if lvText != "" {
			err := c.lv.UnmarshalText([]byte(lvText))
			if err != nil {
				return
			}
		}
	})
	return c
}

func (c *Config) getLevelVariable() string {
	return c.parentKey + ".output.console.level"
}

// SetAutoFormat
// 	@Description  监听 console 格式 的改变
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetAutoFormat() *Config {
	conf.OnChange(func(config *conf.Configuration) {
		formatTxt := config.GetString(c.getFormatVariable())
		if formatTxt != "" {
			c.Format = Format(formatTxt)
		}
	})
	return c
}

func (c *Config) getFormatVariable() string {
	return c.parentKey + ".output.console.format"
}

// Build
// 	@Description 实例化 console 输出源
// 	@Receiver c 配置
// 	@Return zap.Core
func (c Config) Build() zap.Core {
	var ws zap.WriteSyncer
	ws = os.Stdout
	if c.Async {
		var close buffer.CloseFunc
		ws, close = buffer.Buffer(ws, c.BufferSize, ktime.Duration(c.FlushInterval))
		defers.Register(close)
	}

	var encoderCfg EncoderConfig
	if c.Format != "" {
		encoderCfg = *getFormatConsoleEncoderConfig()
		encoderCfg.Format = c.GetFormatValue()
	} else {
		encoderCfg = *getConsoleEncoderConfig()
	}
	if c.lv == nil {
		panic(fmt.Errorf("%s atom level is empty", c.parentKey))
		return nil
	}
	core := zap.NewCore(
		func() zap.Encoder {
			return NewConsoleEncoder(encoderCfg)
		}(),
		ws,
		c.lv,
	)
	return core
}

// getConsoleEncoderConfig 控制台编码配置
func getConsoleEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "logLevel",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zap.DefaultLineEnding,
		EncodeLevel:   zap.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zap.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration:   zap.SecondsDurationEncoder,
		EncodeCaller:     zap.ShortCallerEncoder,
		ConsoleSeparator: " ",
		Format:           property.BuildInExpress(true),
	}
}

// getFormatConsoleEncoderConfig 设置
func getFormatConsoleEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "logLevel",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zap.DefaultLineEnding,
		EncodeLevel:   zap.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zap.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration:   zap.SecondsDurationEncoder,
		EncodeCaller:     zap.ShortCallerEncoder,
		ConsoleSeparator: " ",
		Format:           property.BuildInExpress(true),
	}

}
