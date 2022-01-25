package taskmanager

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

type Configs []TaskConfig

// TaskConfig 配置
type TaskConfig struct {
	// ExecFuncName 业务执行函数名字
	Name string
	// TaskType 任务类型
	TaskType string

	// Required. 触发时间
	//	默认最小单位为分钟.比如:
	//		"* * * * * *" 代表每分钟执行
	//	如果 IsWithSeconds = true. 那么最小单位为秒. 示例:
	//		"*/3 * * * * * *" 代表每三秒钟执行一次
	Spec string
	// IsWithSeconds 使用秒作解析器，默认否
	IsWithSeconds bool
	// IsImmediatelyRun 是否立刻执行，默认否
	IsImmediatelyRun bool

	//IsDistributedTask 是否是分布式任务, 分布式锁
	IsDistributedTask bool
	// DelayExecType skip，queue，concurrent，如果上一个任务执行较慢，到达了新任务执行时间，那么新任务选择跳过，排队，并发执行的策略，新任务默认选择skip策略
	DelayExecType string
}

// RawConfig
// 	@Description tasks 管理配置文件
//	@Param ctx 上下文
//	@Param key 任务管理配置对应的key
func RawConfig(ctx context.Context, key string) Configs {
	var config Configs
	if err := conf.UnmarshalKey(key, &config); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panicf(" key %v unmarshal err:%v", key, err)
	}
	return config
}

// Build
// 	@Description 实例化任务管理
// 	@Receiver c Config
//	@Param ctx 上下文
// 	@Return Manage 任务管理者实例
func (c Configs) Build(ctx context.Context) *Manage {
	mConfig := make(map[string]*TaskConfig, len(c))
	for k := range c {
		cfg := c[k]
		jobName := cfg.Name
		if jobName == "" {
			klog.KuaigoLogger.WithContext(ctx).Panic("task config must have exec function name")
		}
		if _, ok := mConfig[jobName]; ok {
			klog.KuaigoLogger.WithContext(ctx).Panicf("task function name %v repetition in task list", jobName)
		}
		mConfig[jobName] = &cfg
	}
	return &Manage{
		Configs:     mConfig,
		tasksByType: make(map[string][]ktask.Tasker),
		tasksByName: make(map[string]ktask.Tasker),
	}
}
