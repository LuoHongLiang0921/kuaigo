package kuaigo

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

// Construct
//  @Description 构造函数，New完之后必须调用。初始化前置资源
func (app *App) Construct() *App {
	// 初始化函数
	app.initialize()
	return app
}

// GetInstance
//  @Description 获取全局单例，如果new多个，永远指向最后一次Build的实例
func GetInstance() *App {
	instanceOnce.Do(func() {
		new(App).Construct()
	})
	return appInstance
}

// Build
//  @Description 构建服务初始化配置与函数
//  @Receiver app App类型
//  @Param fns 初始化函数集合
func (app *App) Build(fns ...func() error) *App {
	app.RegisterInitFns(fns...)
	return app
}

// Run
//  @Description 服务启动入口，支持启动多个服务
//  @Receiver app App类型
//  @Param servers 服务集合
//  @Return error 启动过程中的错误信息
func (app *App) Run(servers ...server.Server) error {
	// 初始化系统函数
	_ = kgo.SerialUntilError(appInstance.initFns...)()
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	//start signal listen task in goroutine
	app.watchSignals()
	defer app.clean()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	if app.taskManager != nil {
		// todo jobs not graceful
		app.cycle.Run(app.startJobs)

		// start cron task
		app.cycle.Run(app.startCronTasks)
		// start background task
		app.cycle.Run(app.startBackgroundTasks)
	}
	app.runHooks(StageAfterStart)
	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		app.logger.Error("tabby shutdown with error", klog.FieldMod(ecode.ModApp), klog.FieldErr(err))
		return err
	}
	app.logger.Info("shutdown tabby, bye!", klog.FieldMod(ecode.ModApp))
	return nil
}

// Stop
//  @Description 停止服务
//  @Receiver app App类型
//  @Return error 停止过程中的错误信息
func (app *App) Stop() error {
	//var err error
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.RUnlock()
		if app.taskManager != nil {
			//停止定时任务
			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeCron) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.Stop)
				}(w)
			}
			// 停止后台任务
			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeBackground) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.Stop)
				}(w)
			}
			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeOnce) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.Stop)
				}(w)
			}
		}

		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return nil
}

// GracefulStop
//  @Description 优雅停止服务
//  @Receiver app App类型
//  @Param ctx 应用上下文
//  @Return error 停止过程中的错误信息
func (app *App) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)

		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.smu.RUnlock()
		if app.taskManager != nil {
			//stop taskers
			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeCron) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.GracefulStop)
				}(w)
			}
			// 停止后台任务
			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeBackground) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.GracefulStop)
				}(w)
			}

			for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeOnce) {
				func(w ktask.Tasker) {
					app.cycle.Run(w.GracefulStop)
				}(w)
			}

		}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return err
}

// Deprecated: 请使用 WithServiceName 替代
//  @Description 设置服务名字
//  @Receiver app App类型
//  @Param serviceName 服务名称
func (app *App) SetServiceName(serviceName string) {
	app.ServiceName = serviceName
	klog.KuaigoLogger.SetServiceName(app.ServiceName)
}

// WithServiceName
// 	@Description
// 	@Receiver app
//	@Param serviceName
// 	@Return *App
func (app *App) WithServiceName(serviceName string) *App {
	app.ServiceName = serviceName
	klog.KuaigoLogger.SetServiceName(app.ServiceName)
	return app
}

// WithContext
//  @Description 名部设置上下文
//  @Receiver app App类型
//  @Param ctx 上下文实例
//  @Return app 返回本身，方便级联调用
func (app *App) WithContext(ctx context.Context) *App {
	app.ctx = ctx
	return app
}

// GetContext
//  @Description 获取app上下文context
//  @Receiver app App类型
//  @Return context.Context
func (app *App) GetContext() context.Context {
	return app.ctx
}

// RegisterTasks
// 	@Description  添加任务,
// 	@Receiver app App
//	@Param tasks 任务处理函数列表
// 	@Return *App
func (app *App) RegisterTasks(tasks ...ktask.Handler) *App {
	app.smu.Lock()
	defer app.smu.Unlock()
	ctx := app.GetContext()
	ctx = klog.TaskLoggerContext(ctx, klog.WithServiceName(app.ServiceName))
	manager := app.loadOrStoreTaskManager(ctx)
	manager.Load(ctx).RegisterHandlers(ctx, tasks)
	return app
}

// RegisterInitFns
//  @Description 注册初始化函数
//  @Receiver app App类型
//  @Param fns 注册的初始化函数集合
//  @Return app 返回本身，方便级联调用
func (app *App) RegisterInitFns(fns ...func() error) *App {
	app.initFns = append(app.initFns, fns...)
	return app
}

// RegisterServers
//  @Description 注册启动服务
//  @Receiver app App类型
//  @Param fns 注册的启动服务集合
//  @Return app 返回本身，方便级联调用
func (app *App) RegisterServers(s ...server.Server) *App {
	app.servers = append(app.servers, s...)
	return app
}

//RegisterHooks
//  @Description 注册应用钩子
//  @Receiver app App类型
//  @Param fns 注册的应用钩子集合
func (app *App) RegisterHooks(k uint32, fns ...func() error) *App {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Push(fns...)
	}
	return app
}

// Deprecated: 请使用 Build 替代
//  @Description 构建服务初始化函数
//  @Receiver app App类型
//  @Param fns 初始化函数集合
//  @Return error 兼容性返回，永远不会返回报错
func (app *App) Startup(fns ...func() error) error {
	app.Build(fns...)
	return nil
}

// Deprecated: 请使用 RegisterServers 替代
//  @Description 增加启动服务
//  @Receiver app App类型
//  @Param servers 启动服务集合
//  @Return error 兼容性返回，永远不会返回报错
func (app *App) Serve(servers ...server.Server) error {
	app.RegisterServers(servers...)
	return nil
}
