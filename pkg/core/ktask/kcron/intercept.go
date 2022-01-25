package kcron

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
	"time"


	"github.com/robfig/cron/v3"
)

// delayIfStillRunning
// 	@Description 推迟后面的任务调用
//	@Param ctx 上下文
//	@Param logger 日志
// 	@Return JobWrapper 包装后的job
func delayIfStillRunning(ctx context.Context, logger *klog.Logger) JobWrapper {
	return func(j Job) Job {
		var mu sync.Mutex
		return cron.FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.WithContext(ctx).Info("cron delay", klog.String("duration", dur.String()))
			}
			j.Run()
		})
	}
}

// skipIfStillRunning
// 	@Description 如果前一个任务还在运行，则跳过调用
//	@Param ctx
//	@Param logger
// 	@Return JobWrapper
func skipIfStillRunning(ctx context.Context, logger *klog.Logger) JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j Job) Job {
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.WithContext(ctx).Info("cron skip")
			}
		})
	}
}
