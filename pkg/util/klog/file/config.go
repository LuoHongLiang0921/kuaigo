// @Description

package file

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/defers"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/buffer"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/rotate"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"strings"
)

const (
	OutputFile = "file"
)

type Config struct {
	// Level 日志初始等级
	Level string `yaml:"level"`
	// 是否添加调用者信息
	AddCaller  bool `yaml:"addCaller"`
	CallerSkip int  `yaml:"callerSkip"`
	// rootConfigKey output 键名
	rootConfigKey string

	// FileName 文件名
	FileName string `yaml:"fileName"`
	// 日志输出文件最大长度，超过改值则截断
	MaxSize int `yaml:"maxSize"`
	// MaxAge 最大保存时间
	MaxAge int `yaml:"maxAge"`
	// MaxBackup 备份个数
	MaxBackup int `yaml:"maxBackup"`
	// Interval 日志轮换间隔
	Interval string `yaml:"interval"`

	//Async 异步
	Async bool `yaml:"async"`
	// FlushInterval
	FlushInterval string `yaml:"flushInterval"`
	// BufferSize 缓冲大小
	BufferSize int `yaml:"bufferSize"`

	LoggerType string           `mapstructure:"-"`
	parentKey  string           `mapstructure:"-"`
	lv         *zap.AtomicLevel `mapstructure:"-"`
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
		lvText := strings.ToLower(config.GetString(c.parentKey + ".output.file.level"))
		if lvText != "" {
			err := c.lv.UnmarshalText([]byte(lvText))
			if err != nil {
				return
			}
		}
	})
	return c
}

// Build
// 	@Description 构建文件输出源实例
// 	@Receiver c
// 	@Return zap.Core
func (c *Config) Build() zap.Core {
	var ws zap.WriteSyncer
	ws = zap.AddSync(rotate.NewRotateWriter(
		rotate.WithFileName(c.Filename()),
		rotate.WithInterval(ktime.Duration(c.Interval)),
		rotate.WithMaxAge(c.MaxAge),
		rotate.WithMaxSize(c.MaxSize),
		rotate.WithMaxBackups(c.MaxBackup),
	))
	if c.Async {
		var close buffer.CloseFunc
		ws, close = buffer.Buffer(ws, c.BufferSize, ktime.Duration(c.FlushInterval))
		defers.Register(close)
	}
	if c.lv == nil {
		panic(fmt.Errorf("%s atom level is empty", c.parentKey))
		return nil
	}

	fileCore := zap.NewCore(
		func() zap.Encoder {
			return zap.NewJSONEncoder(getFileEncoderConfig())
		}(),
		ws,
		c.lv,
	)
	return fileCore
}

// SetDefaultConfig
// 	@Description 设置默认配置
// 	@Receiver c
// 	@Return *Config
func (c *Config) SetDefaultConfig() *Config {

	if c.FlushInterval == "" {
		c.FlushInterval = "5s"
	}
	if c.BufferSize == 0 {
		c.BufferSize = 256 * 1024
	}
	if c.MaxSize == 0 {
		c.MaxSize = 256 * 1024
	}
	if c.MaxBackup == 0 {
		c.MaxBackup = 10
	}
	if c.MaxAge == 0 {
		c.MaxAge = 1
	}
	if c.Interval == "" {
		c.Interval = "24h"
	}
	return c
}

// Filename
// 	@Description  日志名字
// 	@receiver config
// 	@return string
func (c *Config) Filename() string {
	fileName := c.FileName
	if fileName == "" {
		return c.LoggerType
	}

	leftIndex := strings.Index(fileName, "[")
	if leftIndex < 0 {
		return fileName
	}
	rightIndex := strings.LastIndex(fileName, "]")
	if rightIndex < 0 {
		return fileName
	}
	fileRune := []rune(fileName)

	old := fileRune[leftIndex : rightIndex+1]
	patterStr := fileRune[leftIndex+1 : rightIndex]
	formatStr := ktime.Now().Format(string(patterStr))
	return strings.Replace(fileName, string(old), formatStr, -1)

}

func getFileEncoderConfig() zap.EncoderConfig {
	return zap.EncoderConfig{
		TimeKey:          "timestamp",
		LevelKey:         zap.OmitKey,
		NameKey:          zap.OmitKey,
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stack",
		LineEnding:       zap.DefaultLineEnding,
		EncodeLevel:      zap.LowercaseLevelEncoder,
		EncodeTime:       zap.EpochMillisTimeEncoder,
		EncodeDuration:   zap.SecondsDurationEncoder,
		EncodeCaller:     zap.ShortCallerEncoder,
		EncodeName:       zap.FullNameEncoder,
		ConsoleSeparator: " ",
	}
}
