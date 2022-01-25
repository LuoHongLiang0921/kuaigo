package conf

import (
	"io"

	"github.com/davecgh/go-spew/spew"
)

// ConfigSource ...
type ConfigSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

// Unmarshaller ...
type Unmarshaller = func([]byte, interface{}) error

var defaultConfiguration = New()

// OnChange
//  @Description: 注册change回调函数
//  @Param fn
func OnChange(fn func(*Configuration)) {
	defaultConfiguration.OnChange(fn)
}

// LoadFromConfigSource
//  @Description: 从配置数据源中加载配置
//  @Param ds
//  @Param unmarshaller
//  @Return error
func LoadFromConfigSource(ds ConfigSource, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromConfigSource(ds, unmarshaller)
}

// LoadFromReader
//  @Description:从默认数据配置源中加载配置信息
//  @Param r
//  @Param unmarshaller
//  @Return error
func LoadFromReader(r io.Reader, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromReader(r, unmarshaller)
}

// Apply ...
func Apply(conf map[string]interface{}) error {
	return defaultConfiguration.apply(conf)
}

// Reset
//  @Description: 重置为默认配置
func Reset() {
	defaultConfiguration = New()
}

// Traverse ...
func Traverse(sep string) map[string]interface{} {
	return defaultConfiguration.traverse(sep)
}

// Debug ...
func Debug(sep string) {
	spew.Dump("Debug", Traverse(sep))
}

// Get
//  @Description: returns an interface. For a specific value use one of the Get____ methods.
//  @Param key 配置key值
//  @Return interface{}
func Get(key string) interface{} {
	return defaultConfiguration.Get(key)
}

// Set
//  @Description: 设置配置KV
//  @Param key
//  @Param val
func Set(key string, val interface{}) {
	defaultConfiguration.Set(key, val)
}
