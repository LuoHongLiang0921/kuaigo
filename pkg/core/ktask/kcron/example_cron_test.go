package kcron_test

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

func main() {
	ctx := context.Background()
	eng := NewEngine(ctx)
	if err := eng.Run(); err != nil {
		klog.Fatal(err.Error())
	}
}

type Engine struct {
	*kuaigo.App
}

func NewEngine(ctx context.Context) *Engine {
	eng := new(Engine)
	eng.App = kuaigo.GetInstance()
	eng.Build(
		func() error {
			return eng.initJob(ctx)
		},
	)
	return eng
}
func (e *Engine) initJob(ctx context.Context) error {
	cronTask := NewCronJob("click me ")
	e.WithContext(ctx).RegisterTasks(cronTask)
	return nil
}

type CronJob struct {
	Content string
}

func NewCronJob(content string) *CronJob {
	return &CronJob{Content: content}
}

func (c CronJob) Name() string {
	return "democron"
}

func (c CronJob) BeforeTaskExec(ctx context.Context) error {
	klog.WithContext(ctx).Info("BeforeTaskExec")
	return nil
}

func (c CronJob) AfterTaskExec(ctx context.Context) error {
	klog.WithContext(ctx).Info("AfterTaskExec")
	return nil
}

func (c CronJob) Exec(ctx context.Context, args ...interface{}) error {
	klog.WithContext(ctx).Infof("Exec %v", c.Content)
	return nil
}

func ExampleConfig_Build() {
	ctx := klog.RunningLoggerContext(context.Background())
	eng := NewEngine(ctx)
	if err := eng.Run(); err != nil {
		klog.WithContext(ctx).Error(err.Error())
	}
}
