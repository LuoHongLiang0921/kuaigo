// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package kuaigo

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/biz/kconfig"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/datasource/manager"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/server/governor"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kcycle"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kdefer"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kgo"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/sentinel"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/signals"
	"os"
	"runtime"
	"sync"

	"gopkg.in/yaml.v3"

	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
)

const (
	DisableParserFlag      Disable = 1
	DisableLoadConfig      Disable = 2
	DisableDefaultGovernor Disable = 3
)

const (
	//StageAfterStop after app stop
	StageAfterStop uint32 = iota + 1
	//StageBeforeStop before app stop
	StageBeforeStop
)

type App struct {
	cycle        *kcycle.Cycle
	smu          *sync.RWMutex
	initOnce     sync.Once
	startupOnce  sync.Once
	stopOnce     sync.Once
	modsOnce     sync.Once
	servers      []server.Server
	logger       *klog.Logger
	configParser conf.Unmarshaller
	disableMap   map[Disable]bool
	HideBanner   bool
	hooks        map[uint32]*kdefer.DeferStack
	//ServiceName 上层服务名
	ServiceName string
}

var Instance mods

//mods ...
type mods struct {
	// 配置服务
	*kconfig.AppConfig
}

//Run 启动
// 可以启动多个server
// - 拉取配置
//   - apollo,json
//   - 配置服务(定时拉，缓存)
// - metrics 监控服务
// - pprof 性能监控服务
func (app *App) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	//app.startJobs()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	// start workers
	//app.cycle.Run(app.startWorkers)

	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		app.logger.Error(app.getContext(), "tabby shutdown with error", klog.FieldMod(ecode.ModApp), klog.FieldErr(err))
		return err
	}
	app.logger.Info(app.getContext(), "shutdown tabby, bye!", klog.FieldMod(ecode.ModApp))
	return nil
}

//InitMods ...
// 初始化application
func (app *App) InitMods() error {
	var err error
	app.modsOnce.Do(func() {
		var appConfig kconfig.AppConfig
		err = conf.UnmarshalKey(constant.ConfigAppKey, &appConfig)
		if err != nil {
			return
		}
		if appConfig.ServiceName != "" {
			app.SetServiceName(appConfig.ServiceName)
		}
		Instance.AppConfig = &appConfig
	})
	return err
}

//init hooks
func (app *App) initHooks(hookKeys ...uint32) {
	app.hooks = make(map[uint32]*kdefer.DeferStack, len(hookKeys))
	for _, k := range hookKeys {
		app.hooks[k] = kdefer.NewStack()
	}
}

//run hooks
func (app *App) runHooks(k uint32) {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Clean()
	}
}

//RegisterHooks register a stage Hook
func (app *App) RegisterHooks(k uint32, fns ...func() error) error {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Push(fns...)
		return nil
	}
	return fmt.Errorf("hook stage not found")
}

func (app *App) initialize() {
	app.initOnce.Do(func() {
		klog.KuaigoLogger.SetServiceName(app.ServiceName)
		//assign
		app.cycle = kcycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		//app.workers = make([]worker.Worker, 0)
		//app.jobs = make(map[string]job.Runner)
		app.logger = klog.KuaigoLogger
		app.configParser = yaml.Unmarshal
		app.disableMap = make(map[Disable]bool)
		//private method
		app.initHooks(StageBeforeStop, StageAfterStop)
		//public method
		//app.SetRegistry(registry.Nop{}) //default nop without registry
	})
}

func (app *App) getContext() context.Context {
	return context.TODO()
}

//clean after app quit
func (app *App) clean() {
	_ = klog.RunningLogger.Flush()
	_ = klog.KuaigoLogger.Flush()
	_ = klog.ErrorLogger.Flush()
	_ = klog.AccessLogger.Flush()
	_ = klog.TaskLogger.Flush()
}

//
func (app *App) Stop() error {
	//var err error
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)

		//if app.registerer != nil {
		//	err = app.registerer.Close()
		//	if err != nil {
		//		app.logger.Error("stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		//	}
		//}
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.RUnlock()

		//stop workers
		//for _, w := range app.workers {
		//	func(w worker.Worker) {
		//		app.cycle.Run(w.Stop)
		//	}(w)
		//}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return nil
}

func (app *App) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
	return kgo.SerialUntilError(fns...)()
}

// GracefulStop application after necessary cleanup
func (app *App) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)

		//if app.registerer != nil {
		//	err = app.registerer.Close()
		//	if err != nil {
		//		app.logger.Error("stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		//	}
		//}
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

		////stop workers
		//for _, w := range app.workers {
		//	func(w worker.Worker) {
		//		app.cycle.Run(w.Stop)
		//	}(w)
		//}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return err
}

// waitSignals wait signal
func (app *App) waitSignals() {
	app.logger.Info(app.getContext(), "init listen signal", klog.FieldMod(ecode.ModApp), klog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			app.GracefulStop(context.TODO())
		} else {
			app.Stop()
		}
	})
}
func (app *App) startup() error {
	var err error
	app.startupOnce.Do(func() {
		err = kgo.SerialUntilError(
			app.parseFlags,
			app.printBanner,
			app.loadConfig,
			app.InitMods,
			app.initLogger,
			app.initMaxProcs,
			app.initTracer,
			app.initSentinel,
			app.initGovernor,
		)()
	})
	return err
}

func (app *App) initGovernor() error {
	if app.isDisable(DisableDefaultGovernor) {
		app.logger.Info(app.getContext(), "defualt governor disable", klog.FieldMod(ecode.ModApp))
		return nil
	}

	config := governor.RawConfig("governor")
	if !config.Enable {
		return nil
	}
	return app.Serve(config.Build())
}

func (app *App) startServers() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			//_ = app.registerer.RegisterService(context.TODO(), s.Info())
			//defer app.registerer.UnregisterService(context.TODO(), s.Info())
			app.logger.Info(app.getContext(), "start server", klog.FieldMod(ecode.ModApp), klog.FieldEvent("init"), klog.FieldName(s.Info().Name), klog.FieldAddr(s.Info().Label()), klog.Any("scheme", s.Info().Scheme))
			defer app.logger.Info(app.getContext(), "exit server", klog.FieldMod(ecode.ModApp), klog.FieldEvent("exit"), klog.FieldName(s.Info().Name), klog.FieldErr(err), klog.FieldAddr(s.Info().Label()))
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

func (app *App) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

func (app *App) SetServiceName(s string) {
	app.ServiceName = s
}

//parseFlags init
func (app *App) parseFlags() error {
	if app.isDisable(DisableParserFlag) {
		app.logger.Info(app.getContext(), "parseFlags disable", klog.FieldMod(ecode.ModApp))
		return nil
	}

	flag.Register(&flag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "TABBY_CONFIG",
		Default: "",
		Action:  func(name string, fs *flag.FlagSet) {},
	})

	flag.Register(&flag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: false,
		EnvVar:  "TABBY_CONFIG_WATCH",
	})

	flag.Register(&flag.BoolFlag{
		Name:    "version",
		Usage:   "--version, print version",
		Default: false,
		Action: func(string, *flag.FlagSet) {
			kpkg.PrintVersion()
			os.Exit(0)
		},
	})

	flag.Register(&flag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: "127.0.0.1",
		Action:  func(string, *flag.FlagSet) {},
	})
	return flag.Parse()
}

//loadConfig init
func (app *App) loadConfig() error {
	if app.isDisable(DisableLoadConfig) {
		app.logger.Info(app.getContext(), "load config disable", klog.FieldMod(ecode.ModConfig))
		return nil
	}

	var configAddr = flag.String("config")
	if configAddr == "" {
		configAddr = os.Getenv("APOLLO_SERVER_ADDR")
	}
	provider, err := manager.NewDataSource(configAddr)
	if err != manager.ErrConfigAddr {
		if err != nil {
			app.logger.Panic(app.getContext(), "data source: provider error", klog.FieldMod(ecode.ModConfig), klog.FieldErr(err))
		}

		if err := conf.LoadFromDataSource(provider, app.configParser); err != nil {
			app.logger.Panic(app.getContext(), "data source: load config", klog.FieldMod(ecode.ModConfig), klog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), klog.FieldErr(err))
		}
	} else {
		app.logger.Info(app.getContext(), "no config... ", klog.FieldMod(ecode.ModConfig))
	}
	return nil
}

//initLogger init
// 从apollo中拉取
func (app *App) initLogger() error {
	if conf.Get("running") != nil {
		klog.RunningLogger = klog.RawConfig("running").Build()
	}
	klog.RunningLogger.AutoLevel("running")
	klog.RunningLogger.SetServiceName(app.ServiceName)
	if conf.Get("tabbyLogger") != nil {
		klog.KuaigoLogger = klog.RawConfig("kuaigoLogger").Build()
	}
	klog.KuaigoLogger.SetServiceName(app.ServiceName)
	klog.KuaigoLogger.AutoLevel("tabbyLogger")
	if conf.Get("access") != nil {
		klog.AccessLogger = klog.RawConfig("access").Build()
	}
	klog.AccessLogger.SetServiceName(app.ServiceName)
	klog.AccessLogger.AutoLevel("access")
	if conf.Get("error") != nil {
		klog.ErrorLogger = klog.RawConfig("error").Build()
	}
	klog.ErrorLogger.SetServiceName(app.ServiceName)
	klog.ErrorLogger.AutoLevel("error")
	if conf.Get("task") != nil {
		klog.TaskLogger = klog.RawConfig("task").Build()
	}
	klog.TaskLogger.SetServiceName(app.ServiceName)
	klog.TaskLogger.AutoLevel("task")
	return nil
}

//initTracer init
func (app *App) initTracer() error {
	// init tracing component jaeger
	//if conf.Get("tabby.trace.jaeger") != nil {
	//	var config = jaeger.RawConfig("tabby.trace.jaeger")
	//	trace.SetGlobalTracer(config.Build())
	//}
	return nil
}

//initSentinel init
func (app *App) initSentinel() error {
	// init reliability component sentinel
	if conf.Get("tabby.reliability.sentinel") != nil {
		app.logger.Info(app.getContext(), "init sentinel")
		return sentinel.RawConfig("tabby.reliability.sentinel").Build()
	}
	return nil
}

//initMaxProcs init
func (app *App) initMaxProcs() error {
	if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			app.logger.Panic(app.getContext(), "auto max procs", klog.FieldMod(ecode.ModProc), klog.FieldErrKind(ecode.ErrKindAny), klog.FieldErr(err))
		}
	}
	app.logger.Info(app.getContext(), "auto max procs", klog.FieldMod(ecode.ModProc), klog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	return nil
}

func (app *App) isDisable(d Disable) bool {
	b, ok := app.disableMap[d]
	if !ok {
		return false
	}
	return b
}

//printBanner init
func (app *App) printBanner() error {
	if app.HideBanner {
		return nil
	}
	const banner = `Welcome to tabby, starting application ...`
	fmt.Println(kcolor.Green(banner))
	return nil
}
