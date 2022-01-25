// @Description 流量控制，熔断降级

package sentinel

import (
	"context"
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"io/ioutil"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	sentinelConfig "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("tabby.sentinel." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		klog.Panic("unmarshal key", klog.Any("err", err))
	}
	return config
}

// Config ...
type Config struct {
	AppName       string       `json:"appName"`
	LogPath       string       `json:"logPath"`
	FlowRules     []*flow.Rule `json:"flowRules"`
	FlowRulesFile string       `json:"flowRulesFile"`
}

// DefaultConfig returns default config for sentinel
func DefaultConfig() *Config {
	return &Config{
		AppName:   pkg.GetAppName(),
		LogPath:   "/tmp/log",
		FlowRules: make([]*flow.Rule, 0),
	}
}
func (config *Config) getContext() context.Context {
	return context.TODO()
}

// Build InitSentinelCoreComponent init sentinel core component
// Currently, only flow rules from json file is supported
// todo: support dynamic rule config
// todo: support more rule such as system rule
func (config *Config) Build() error {
	if config.FlowRulesFile != "" {
		var rules []*flow.Rule
		content, err := ioutil.ReadFile(config.FlowRulesFile)
		if err != nil {
			klog.Error("load sentinel flow rules", klog.FieldErr(err), klog.FieldKey(config.FlowRulesFile))
		}

		if err := json.Unmarshal(content, &rules); err != nil {
			klog.Error("load sentinel flow rules", klog.FieldErr(err), klog.FieldKey(config.FlowRulesFile))
		}

		config.FlowRules = append(config.FlowRules, rules...)
	}

	configEntity := sentinelConfig.NewDefaultConfig()
	configEntity.Sentinel.App.Name = config.AppName
	configEntity.Sentinel.Log.Dir = config.LogPath

	if len(config.FlowRules) > 0 {
		_, _ = flow.LoadRules(config.FlowRules)
	}
	return sentinel.InitWithConfig(configEntity)
}

func Entry(resource string) (*base.SentinelEntry, *base.BlockError) {
	return sentinel.Entry(resource)
}
