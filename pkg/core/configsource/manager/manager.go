package manager

import (
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"net/url"

)

var (
	//ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidConfigSource defines an error that the scheme has been registered
	ErrInvalidConfigSource = errors.New("invalid configsource, please make sure the scheme has been registered")
	registry               map[string]ConfigSourceCreatorFunc
	//DefaultScheme ..
	DefaultScheme string
)

// ConfigSourceCreatorFunc represents a configSource creator function
type ConfigSourceCreatorFunc func() conf.ConfigSource

func init() {
	registry = make(map[string]ConfigSourceCreatorFunc)
}

// Register registers a configSource creator function to the registry
func Register(scheme string, creator ConfigSourceCreatorFunc) {
	registry[scheme] = creator
}

// NewConfigSource
// 	@Description
//	@param configAddr 获取schema对应的ConfigSource，并且执行
// 	@return conf.ConfigSource
// 	@return error
func NewConfigSource(configAddr string) (conf.ConfigSource, error) {
	if configAddr == "" {
		return nil, ErrConfigAddr
	}
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		DefaultScheme = urlObj.Scheme
	}

	creatorFunc, exist := registry[DefaultScheme]
	if !exist {
		return nil, ErrInvalidConfigSource
	}
	return creatorFunc(), nil
}
