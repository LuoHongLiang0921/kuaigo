package kcron

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Config 配置
type Config struct {
	Name string
	// Required. 触发时间
	//	默认最小单位为分钟.比如:
	//		"* * * * * *" 代表每分钟执行
	//	如果 EnableSeconds = true. 那么最小单位为秒. 示例:
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
	wrappers      []JobWrapper
	logger        *klog.Logger
	parser        cron.Parser
}

// RawConfig
// 	@Description 实例配置
//	@Param ctx 上下文
//	@Param key
// 	@Return Config 实例后的
func RawConfig(ctx context.Context, key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panicf("key %v unmarshal cron RawConfig", key)
	}
	config.logger = klog.KuaigoLogger
	if config.IsDistributedTask {
		// todo:待实现
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		logger:           klog.KuaigoLogger,
		wrappers:         []JobWrapper{},
		IsWithSeconds:    false,
		IsImmediatelyRun: false,
		DelayExecType:    "skip", // skip
	}
}

// Build
// 	@Description 实例化定时任务
// 	@Receiver config 配置
//	@Param ctx 上下文
// 	@Return *XCron
func (config Config) Build(ctx context.Context) *XCron {
	if config.IsWithSeconds {
		config.parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	} else {
		// default parser
		config.parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	switch config.DelayExecType {
	case "skip":
		config.wrappers = append(config.wrappers, skipIfStillRunning(ctx, config.logger))
	case "queue":
		config.wrappers = append(config.wrappers, delayIfStillRunning(ctx, config.logger))
	case "concurrent":
	default:
		config.wrappers = append(config.wrappers, skipIfStillRunning(ctx, config.logger))
	}

	if config.IsDistributedTask {

	}

	return newCron(&config)
}

// WithLogger
// 	@Description 设置logger
// 	@Receiver config 配置
//	@Param lg 日志实例
// 	@Return *Config 设置日志实例后的logger
func (config *Config) WithLogger(lg *klog.Logger) *Config {
	config.logger = lg
	return config
}

type wrappedJob struct {
	NamedJob
	logger            *klog.Logger
	IsDistributedTask bool
	ctx               context.Context
	debug             bool
}

// Run
// 	@Description  运行job
// 	@Receiver wj wrappedJob
func (wj wrappedJob) Run() {
	if wj.IsDistributedTask {
		//Todo: 获取锁
	}
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	var fields = []klog.Field{zap.String("name", wj.Name())}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, klog.String("err", err.Error()), klog.Duration("cost", time.Since(beg)))
			wj.logger.WithContext(wj.ctx).Error("run", fields...)
		} else {
			if wj.debug {
				wj.logger.WithContext(wj.ctx).Infof("run job %v", wj.Name())
			}
		}
	}()

	return wj.NamedJob.Run()
}
