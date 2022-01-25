// @Description

package klog

import "fmt"

const (
	// flagFile 文件 1
	flagFile = 1
	// flagRedis redis 2
	flagRedis = 1 << 1
	// flagConsole console 4
	flagConsole = 1 << 2

	StoreFile    = "file"
	StoreRedis   = "redis"
	StoreConsole = "console"
)

var (
	storeMap = map[string]int{
		"file":    flagFile,
		"redis":   flagRedis,
		"console": flagConsole,
	}
)

func (c *Config) getStore() (int, error) {
	store := flagConsole
	for _, v := range c.Store {
		if mode, ok := storeMap[v]; ok {
			if mode == store {
				continue
			}
			store = store | mode
		} else {
			return 0, fmt.Errorf("not suporrt %v", c.Store)
		}
	}
	return store, nil
}

// SetStore
// 	@Description  设置存储状态
// 	@receiver c
//	@Param status
func (c *Config) setStore(status int) *Config {
	c.store = status
	return c
}

// AddStore
// 	@Description  增加存储状态
// 	@receiver c
//	@Param s
func (c *Config) addStore(s int) *Config {
	c.store = c.store | s
	return c
}

// DeleteStore
// 	@Description  删除状态
// 	@receiver c
//	@Param s
func (c *Config) deleteStore(s int) *Config {
	c.store &= ^s
	return c
}

// HasStore
// 	@Description  存储状态是否包含 s 状态
// 	@receiver c
//	@Param s
// 	@return bool
func (c *Config) hasStore(s int) bool {
	return (c.store & s) == s
}

// NotHasStore
// 	@Description  是否不具有某些状态
// 	@receiver c
//	@Param s
// 	@return bool
func (c *Config) notHasStore(s int) bool {
	return (c.store & s) == 0
}

// OnlyHasStore
// 	@Description  仅包含某种状态
// 	@receiver c
//	@Param s
// 	@return bool
func (c *Config) onlyHasStore(s int) bool {
	return c.store == s
}
