package governor

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/knet"
)

//ModName ..
const ModName = "govern"

// Config ...
type Config struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Network string `json:"network" toml:"network" yaml:"network"`
	logger  *klog.Logger
	Enable  bool `yaml:"enable"`

	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string
}

// StdConfig represents Standard gRPC Server config
// which will parse config by conf package,
// panic if no config key found in conf
func StdConfig(name string) *Config {
	return RawConfig("tabby.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Get(key) == nil {
		return config
	}
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("govern server parse config panic",
			klog.FieldErr(err), klog.FieldKey(key),
			klog.FieldValueAny(config),
		)
	}
	return config
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	host, port, err := knet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}

	return &Config{
		Enable:  true,
		Host:    host,
		Network: "tcp4",
		Port:    port,
		logger:  klog.KuaigoLogger.With(klog.FieldMod(ModName)),
	}
}

// Build ...
func (config *Config) Build() *Server {
	return newServer(config)
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
