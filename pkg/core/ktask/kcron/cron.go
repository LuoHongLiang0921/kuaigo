package kcron

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kstring"

	"github.com/robfig/cron/v3"
)

var (
	// Every ...
	Every = cron.Every
	// NewParser ...
	NewParser = cron.NewParser
	// NewChain ...
	NewChain = cron.NewChain
	// WithSeconds ...
	WithSeconds = cron.WithSeconds
	// WithParser ...
	WithParser = cron.WithParser
	// WithLocation ...
	WithLocation = cron.WithLocation
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Entry ...
	Entry = cron.Entry
	// Schedule ...
	Schedule = cron.Schedule
	// Parser ...
	Parser = cron.Parser
	// Option ...
	Option = cron.Option
	// Job ...
	Job = cron.Job
	//NamedJob ..
	NamedJob interface {
		Run() error
		Name() string
	}
)

// FuncJob ...
type FuncJob func() error

func (f FuncJob) Run() error { return f() }

func (f FuncJob) Name() string { return kstring.FunctionName(f) }

// XCron ...
type XCron struct {
	Config  *Config
	Cron    *cron.Cron
	entries map[string]EntryID
	ctx     context.Context
	handler ktask.Handler
	Debug   bool
}

func newCron(config *Config) *XCron {
	if config.logger == nil {
		config.logger = klog.KuaigoLogger
	}
	cron := &XCron{
		Config: config,
		Cron: cron.New(
			cron.WithParser(config.parser),
			cron.WithChain(config.wrappers...),
		),
	}
	return cron
}

// schedule
// 	@Description 添加任务到调度列表
// 	@Receiver c XCron
//	@Param ctx 上下文
//	@Param schedule 调度策略
//	@Param job 业务
// 	@Return EntryID
func (c *XCron) schedule(ctx context.Context, schedule Schedule, job NamedJob) EntryID {
	if c.Config.IsImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innerJob := &wrappedJob{
		NamedJob:          job,
		logger:            c.Config.logger,
		debug:             c.Debug,
		IsDistributedTask: c.Config.IsDistributedTask,
	}
	if c.Debug {
		c.Config.logger.WithContext(ctx).Infof("add job name %v", job.Name())
	}
	return c.Cron.Schedule(schedule, innerJob)
}

// GetEntryByName
// 	@Description
// 	@Receiver c
//	@Param name
// 	@Return cron.Entry
func (c *XCron) GetEntryByName(name string) cron.Entry {
	// todo(gorexlv): data race
	return c.Cron.Entry(c.entries[name])
}

// addJob
// 	@Description 添加工作
// 	@Receiver c XCron
//	@Param ctx 上下文
//	@Param spec 触发规则
//	@Param cmd 处理函数
// 	@Return EntryID 业务job
// 	@Return error 错误
func (c *XCron) addJob(ctx context.Context, cmd NamedJob) (EntryID, error) {
	schedule, err := c.Config.parser.Parse(c.Config.Spec)
	if err != nil {
		return 0, err
	}
	return c.schedule(ctx, schedule, cmd), nil
}

// RegisterHandler
// 	@Description 注册业务处理函数
// 	@Receiver c XCron
//	@Param ctx 上下文
//	@Param spec 规则
//	@Param handler 处理函数
// 	@Return EntryID
// 	@Return error
func (c *XCron) RegisterHandler(ctx context.Context, handler ktask.Handler) error {
	c.handler = handler
	if c.handler == nil {
		return fmt.Errorf("task name %v not register handler is nil", c.Name())
	}
	_, err := c.addJob(ctx, FuncJob(func() error {
		return c.handler.Exec(ctx)
	}))
	return err
}

// Name
// 	@Description 任务名字
// 	@Receiver c XCron
// 	@Return string
func (c *XCron) Name() string {
	return c.Config.Name
}

// TaskType
// 	@Description 任务类型
// 	@Receiver c XCron
// 	@Return string cron 类型
func (c *XCron) TaskType() string {
	return constant.TaskTypeCron
}

// Run
// 	@Description 启动任务
// 	@Receiver c XCron
// 	@Return error
func (c *XCron) Run() error {
	if c.Debug {
		c.Config.logger.WithContext(c.ctx).Infof("run xtask number of scheduled jobs %v", len(c.Cron.Entries()))
	}
	if c.handler == nil {
		return fmt.Errorf("task name %v not register handler ", c.Name())
	}
	c.handler.BeforeTaskExec(c.ctx)
	c.Cron.Run()
	c.handler.AfterTaskExec(c.ctx)
	return nil
}

// Stop
// 	@Description  停止任务
// 	@Receiver c XCron
// 	@Return error 错误
func (c *XCron) Stop() error {
	_ = c.Cron.Stop()
	return nil
}

func (c *XCron) GracefulStop() error {
	// 检测是否还有消息
	// 有等待，加超时
	_ = c.Cron.Stop()
	return nil
}
