package conf

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcast"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kmap"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Configuration provides configuration for application.
type Configuration struct {
	mu       sync.RWMutex
	override map[string]interface{}
	keyDelim string

	keyMap    *sync.Map
	onChanges []func(*Configuration)

	watchers map[string][]func(*Configuration)
}

const (
	defaultKeyDelim = "."
)

// New
//  @Description:实例化配置信息
//  @Return *Configuration
func New() *Configuration {
	return &Configuration{
		override:  make(map[string]interface{}),
		keyDelim:  defaultKeyDelim,
		keyMap:    &sync.Map{},
		onChanges: make([]func(*Configuration), 0),
		watchers:  make(map[string][]func(*Configuration)),
	}
}

// SetKeyDelim
//  @Description: 为配置项设置key分隔符
//  @Receiver c
//  @Param delim
func (c *Configuration) SetKeyDelim(delim string) {
	c.keyDelim = delim
}

// Sub returns new Configuration instance representing a sub tree of this instance.
func (c *Configuration) Sub(key string) *Configuration {
	return &Configuration{
		keyDelim: c.keyDelim,
		override: c.GetStringMap(key),
	}
}

// WriteConfig
//  @Description: 写配置
//  @Receiver c
//  @Return error
func (c *Configuration) WriteConfig() error {
	//return c.provider.Write(c.override)
	return nil
}

// OnChange
//  @Description: 注册change回调函数
//  @Receiver c
//  @Param fn
func (c *Configuration) OnChange(fn func(*Configuration)) {
	c.onChanges = append(c.onChanges, fn)
}

// LoadFromConfigSource
//  @Description  从配置数据源中加载配置
//  @Receiver c
//  @Param ds
//  @Param unmarshaller
//  @Return error
func (c *Configuration) LoadFromConfigSource(ds ConfigSource, unmarshaller Unmarshaller) error {
	content, err := ds.ReadConfig()
	if err != nil {
		return err
	}

	if err := c.Load(content, unmarshaller); err != nil {
		return err
	}

	go func() {
		for range ds.IsConfigChanged() {
			if content, err := ds.ReadConfig(); err == nil {
				_ = c.Load(content, unmarshaller)
				for _, change := range c.onChanges {
					change(c)
				}
			}
		}
	}()

	return nil
}

// Load
//  @Description  load配置
//  @Receiver c
//  @Param content
//  @Param unmarshal
//  @Return error
func (c *Configuration) Load(content []byte, unmarshal Unmarshaller) error {
	configuration := make(map[string]interface{})
	if err := unmarshal(content, &configuration); err != nil {
		return err
	}
	return c.apply(configuration)
}

// LoadFromReader
//  @Description: 从配置源中加载配置信息
//  @Receiver c
//  @Param reader
//  @Param unmarshaller
//  @Return error
func (c *Configuration) LoadFromReader(reader io.Reader, unmarshaller Unmarshaller) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return c.Load(content, unmarshaller)
}

func (c *Configuration) apply(conf map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var changes = make(map[string]interface{})

	kmap.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		orig, ok := c.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		c.keyMap.Store(k, v)
	}

	if len(changes) > 0 {
		c.notifyChanges(changes)
	}

	return nil
}

// notifyChanges
//  @Description: 配置变更通知
//  @Receiver c
//  @Param changes
func (c *Configuration) notifyChanges(changes map[string]interface{}) {
	var changedWatchPrefixMap = map[string]struct{}{}

	for watchPrefix := range c.watchers {
		for key := range changes {
			// 前缀匹配即可
			// todo 可能产生错误匹配
			if strings.HasPrefix(key, watchPrefix) {
				changedWatchPrefixMap[watchPrefix] = struct{}{}
			}
		}
	}

	for changedWatchPrefix := range changedWatchPrefixMap {
		for _, handle := range c.watchers[changedWatchPrefix] {
			go handle(c)
		}
	}
}

// Set
//  @Description: set配置信息
//  @Receiver c
//  @Param key
//  @Param val
//  @Return error
func (c *Configuration) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
	// c.keyMap.Store(key, val)
}

// deepSearch
//  @Description  deepSearch遍历map
//  @Param m
//  @Param path
//  @Return map[string]interface{}
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}

// Get
//  @Description returns the value associated with the key
//  @Receiver c
//  @Param key
//  @Return interface{}
func (c *Configuration) Get(key string) interface{} {
	return c.find(key)
}

// GetString
//  @Description  字符串形式返回value
//  @Param key
//  @Return string
func GetString(key string) string {
	return defaultConfiguration.GetString(key)
}

// GetString returns the value associated with the key as a string.
func (c *Configuration) GetString(key string) string {
	return kcast.ToString(c.Get(key))
}

// GetBool
//  @Description  返回布尔值
//  @Param key
//  @Return bool
func GetBool(key string) bool {
	return defaultConfiguration.GetBool(key)
}

// GetBool returns the value associated with the key as a boolean.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return bool
func (c *Configuration) GetBool(key string) bool {
	return kcast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return int
func GetInt(key string) int {
	return defaultConfiguration.GetInt(key)
}

// GetInt returns the value associated with the key as an integer.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return int
func (c *Configuration) GetInt(key string) int {
	return kcast.ToInt(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return int64
func GetInt64(key string) int64 {
	return defaultConfiguration.GetInt64(key)
}

// GetInt64 returns the value associated with the key as an integer.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return int64
func (c *Configuration) GetInt64(key string) int64 {
	return kcast.ToInt64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64 with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return float64
func GetFloat64(key string) float64 {
	return defaultConfiguration.GetFloat64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return float64
func (c *Configuration) GetFloat64(key string) float64 {
	return kcast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return time.Time
func GetTime(key string) time.Time {
	return defaultConfiguration.GetTime(key)
}

// GetTime returns the value associated with the key as time.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return time.Time
func (c *Configuration) GetTime(key string) time.Time {
	return kcast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration with default defaultConfiguration.
func GetDuration(key string) time.Duration {
	return defaultConfiguration.GetDuration(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *Configuration) GetDuration(key string) time.Duration {
	return kcast.ToDuration(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return []string
func GetStringSlice(key string) []string {
	return defaultConfiguration.GetStringSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return []string
func (c *Configuration) GetStringSlice(key string) []string {
	return kcast.ToStringSlice(c.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return []interface{}
func GetSlice(key string) []interface{} {
	return defaultConfiguration.GetSlice(key)
}

// GetSlice returns the value associated with the key as a slice of strings.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return []interface{}
func (c *Configuration) GetSlice(key string) []interface{} {
	return kcast.ToSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return map[string]interface{}
func GetStringMap(key string) map[string]interface{} {
	return defaultConfiguration.GetStringMap(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return map[string]interface{}
func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return kcast.ToStringMap(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return map[string]string
func GetStringMapString(key string) map[string]string {
	return defaultConfiguration.GetStringMapString(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return map[string]string
func (c *Configuration) GetStringMapString(key string) map[string]string {
	return kcast.ToStringMapString(c.Get(key))
}

// GetSliceStringMap returns the value associated with the slice of maps.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return []map[string]interface{}
func (c *Configuration) GetSliceStringMap(key string) []map[string]interface{} {
	return kcast.ToSliceStringMap(c.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings with default defaultConfiguration.
// 	@Description
//	@Param key
// 	@Return map[string][]string
func GetStringMapStringSlice(key string) map[string][]string {
	return defaultConfiguration.GetStringMapStringSlice(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
// 	@Description
// 	@Receiver c
//	@Param key
// 	@Return map[string][]string
func (c *Configuration) GetStringMapStringSlice(key string) map[string][]string {
	return kcast.ToStringMapStringSlice(c.Get(key))
}

// UnmarshalWithExpect unmarshal key, returns expect if failed
// 	@Description
//	@Param key
//	@Param expect
// 	@Return interface{}
func UnmarshalWithExpect(key string, expect interface{}) interface{} {
	return defaultConfiguration.UnmarshalWithExpect(key, expect)
}

// UnmarshalWithExpect unmarshal key, returns expect if failed
// 	@Description
// 	@Receiver c
//	@Param key
//	@Param expect
// 	@Return interface{}
func (c *Configuration) UnmarshalWithExpect(key string, expect interface{}) interface{} {
	err := c.UnmarshalKey(key, expect)
	if err != nil {
		return expect
	}
	return expect
}

// UnmarshalKey 接受一个键并将其解组为具有默认 defaultConfiguration 的 Struct。
// 	@Description
//	@Param key
//	@Param rawVal
//	@Param opts
// 	@Return error
func UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	return defaultConfiguration.UnmarshalKey(key, rawVal, opts...)
}

// ErrInvalidKey ...
var ErrInvalidKey = errors.New("invalid key, maybe not exist in config")

// UnmarshalKey 获取一个键并将其解组为一个结构体。
// 	@Description
// 	@Receiver c
//	@Param key
//	@Param rawVal
//	@Param opts
// 	@Return error
func (c *Configuration) UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	var options = defaultGetOptions
	for _, opt := range opts {
		opt(&options)
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    options.TagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.mu.RLock()
		defer c.mu.RUnlock()
		return decoder.Decode(c.override)
	}

	value := c.Get(key)
	if value == nil {
		return errors.Wrap(ErrInvalidKey, key)
	}

	return decoder.Decode(value)
}

func (c *Configuration) find(key string) interface{} {
	dd, ok := c.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, c.keyDelim)
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := kmap.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.keyMap.Store(key, dd)
	return dd
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = k
		}
		if dd, err := kcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}

func (c *Configuration) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", c.override, data, sep)
	return data
}
