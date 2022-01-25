package kuaigo

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/apollo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/file"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/http"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask/taskmanager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/net/kthrift"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/governor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/defers"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/pkg/sentinel"
	"github.com/LuoHongLiang0921/kuaigo/pkg/signals"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcycle"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kdefer"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"os"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"go.uber.org/automaxprocs/maxprocs"
)

// startServers
//  @Description 真正开始启动服务
//  @Receiver app App类型
//  @Return error 启动服务过程中的报错
func (app *App) startServers() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			app.logger.Info("start server", klog.FieldMod(ecode.ModApp), klog.FieldEvent("init"), klog.FieldName(s.Info().Name), klog.FieldAddr(s.Info().Label()), klog.Any("scheme", s.Info().Scheme))
			defer app.logger.Info("exit server", klog.FieldMod(ecode.ModApp), klog.FieldEvent("exit"), klog.FieldName(s.Info().Name), klog.FieldErr(err), klog.FieldAddr(s.Info().Label()))
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

// clean
//  @Description 服务退出时清理资源
//  @Receiver app App类型
//  @Return error 启动服务过程中的报错
func (app *App) clean() {
	klog.FlushAll()
	_ = kthrift.CloseService(app.GetContext())
	defers.Execute()
}

// watchSignals
//  @Description  检测信号
//  @Receiver app App类型
func (app *App) watchSignals() {
	app.logger.Info("init listen signal", klog.FieldMod(ecode.ModApp), klog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			_ = app.GracefulStop(context.TODO())
		} else {
			_ = app.Stop()
		}
	})
}

// runHooks
//  @Description 服务退出时运行钩子
//  @Receiver app App类型
func (app *App) runHooks(k uint32) {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Execute()
	}
}

// initialize
//  @Description app构建后初始化的无状态资源
//  @Receiver app App类型
func (app *App) initialize() {
	app.buildOnce.Do(func() {
		appInstance = app
		klog.KuaigoLogger.SetServiceName(app.ServiceName)
		//assign
		app.cycle = kcycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.initFns = []func() error{
			app.printBanner,
			app.parseFlags,
			app.loadConfig,
			app.initConfigVersion,
			app.initLogger,
			app.initMaxProcess,
			app.initSentinel,
			app.initGovernor,
		}
		app.servers = make([]server.Server, 0)
		app.hooks = map[uint32]*kdefer.DeferStack{
			StageAfterStart: kdefer.NewStack(),
			StageBeforeStop: kdefer.NewStack(),
			StageAfterStop:  kdefer.NewStack(),
		}
		app.logger = klog.KuaigoLogger
		app.configParser = yaml.Unmarshal
		app.disableMap = make(map[Disable]bool)
	})
}

// startTasks
//  @Description 启动定时任务
//  @Receiver app App类型
//  @Return error 任务启动过程中的报错
func (app *App) startCronTasks() error {
	var eg errgroup.Group
	for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeCron) {
		w := w
		eg.Go(func() error {
			return w.Run()
		})
	}
	return eg.Wait()
}

// startBackgroundTasks
// 	@Description 启动后台任务
// 	@Receiver app App
// 	@Return error 任务启动过程中错误
func (app *App) startBackgroundTasks() error {
	var eg errgroup.Group
	for _, w := range app.taskManager.GetTasksByType(constant.TaskTypeBackground) {
		w := w
		eg.Go(func() error {
			return w.Run()
		})
	}
	return eg.Wait()
}

// startJobs
//  @Description 启动一次性Job任务
//  @Receiver app App类型
//  @Return error 任务启动过程中的报错
func (app *App) startJobs() error {
	if flag.Bool("disable-once-job") {
		app.logger.Info("tabby disable once job")
		return nil
	}

	jobFlag := flag.String("job")
	var tasks []ktask.Tasker
	if jobFlag == "" {
		app.logger.Warn("tabby jobs flag name empty,run all once job!")
		tasks = app.taskManager.GetTasksByType(constant.TaskTypeOnce)
	} else {
		taskNames := strings.Split(jobFlag, ",")
		for _, v := range taskNames {
			onceTask := app.taskManager.GetTaskByName(v)
			if onceTask == nil {
				app.logger.Warn("tabby jobs flag name not in tasks,please check config file")
				continue
			}
			tasks = append(tasks, onceTask)
		}
	}

	if len(tasks) <= 0 {
		return nil
	}

	var eg errgroup.Group
	for _, w := range tasks {
		w := w
		eg.Go(func() error {
			app.logger.Infof("job run begin %v", w.Name())
			defer app.logger.Infof("job run end %v", w.Name())
			return w.Run()
		})
	}
	return eg.Wait()
}

// parseFlags
//  @Description 解析命令行参数
//  @Receiver app App类型
// 	@Return error 解析参数时的报错
func (app *App) parseFlags() error {
	if app.isDisable(DisableParserFlag) {
		app.logger.Info("parseFlags disable", klog.FieldMod(ecode.ModApp))
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
			pkg.PrintVersion()
			os.Exit(0)
		},
	})

	flag.Register(&flag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: "127.0.0.1",
		Action:  func(string, *flag.FlagSet) {},
	})

	flag.Register(
		&flag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
	flag.Register(
		&flag.BoolFlag{
			Name:    "disable-once-job",
			Usage:   "--disable-once-job",
			Default: true,
		},
	)
	return flag.Parse()
}

// printBanner
//  @Description 打印启动Banner
//  @Receiver app App类型
// 	@Return error 解析参数时的报错
func (app *App) printBanner() error {
	if app.HideBanner {
		return nil
	}
	fmt.Println(kcolor.Green(string(constant.LogoContent)))
	const banner = `Welcome to kuaigo, starting application ...`
	fmt.Println(kcolor.Green(banner))
	return nil
}

// loadConfig
//  @Description 加载服务配置项
//  @Receiver app App类型
// 	@Return error 加载服务配置时的报错
func (app *App) loadConfig() error {
	if app.isDisable(DisableLoadConfig) {
		app.logger.Info("load config disable", klog.FieldMod(ecode.ModConfig))
		return nil
	}

	file.RegisterConfigHandler()
	http.RegisterConfigHandler()
	apollo.RegisterConfigHandler()
	var configAddr = flag.String("config")
	if configAddr == "" {
		configAddr = os.Getenv("APOLLO_SERVER_ADDR")
	}
	provider, err := manager.NewConfigSource(configAddr)
	if err != manager.ErrConfigAddr {
		if err != nil {
			app.logger.Panic("data source: provider error", klog.FieldMod(ecode.ModConfig), klog.FieldErr(err))
		}

		if err := conf.LoadFromConfigSource(provider, app.configParser); err != nil {
			app.logger.Panic("data source: load config", klog.FieldMod(ecode.ModConfig), klog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), klog.FieldErr(err))
		}
	} else {
		app.logger.Panic("Missing config... ", klog.FieldMod(ecode.ModConfig))
	}
	return nil
}

// initConfigVersion
//  @Description 初始化版本信息
//  @Receiver app App类型
// 	@Return error 初始化版本信息时的报错
func (app *App) initConfigVersion() error {
	var v string
	_ = conf.UnmarshalKey("version", &v)
	app.ConfigVersion = v
	if app.ConfigVersion > "" {
		fmt.Printf("%s:%s/\r\n", kcolor.Green("pool WithConfigVersion:"), app.ConfigVersion)
	} else {
		fmt.Printf("%s:old\r\n", kcolor.Green("pool WithConfigVersion:"))
	}
	return nil
}

// initLogger
//  @Description 初始化日志配置（兼容新版日志配置结构）
//  @Receiver app App类型
// 	@Return error 初始化日志时的报错
func (app *App) initLogger() error {
	klog.InitLogger(app.ConfigVersion, app.ServiceName)
	return nil
}

// initMaxProcess
//  @Description 初始化最大使用Logical Processor cpu核数
//  @Receiver app App类型
// 	@Return error 初始化最大进程数的报错
func (app *App) initMaxProcess() error {
	if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			app.logger.Panic("auto max procs", klog.FieldMod(ecode.ModProc), klog.FieldErrKind(ecode.ErrKindAny), klog.FieldErr(err))
		}
	}
	fmt.Printf("%s:%v\r\n", kcolor.Green("init Max Procs:"), int64(runtime.GOMAXPROCS(-1)))
	return nil
}

// initSentinel
//  @Description 初始化Sentinel限流
//  @Receiver app App类型
// 	@Return error 初始化限流的报错
func (app *App) initSentinel() error {
	// init reliability component sentinel
	//todo: start tabby ?,constant
	if conf.Get("tabby.reliability.sentinel") != nil {
		app.logger.Info("init sentinel")
		return sentinel.RawConfig("tabby.reliability.sentinel").Build()
	}
	return nil
}

// initGovernor
//  @Description 初始化性能，服务信息、指标收集，pprof 信息
//  @Receiver app App类型
// 	@Return error 初始化时报错
func (app *App) initGovernor() error {
	if app.isDisable(DisableDefaultGovernor) {
		app.logger.Info("default governor disable", klog.FieldMod(ecode.ModApp))
		return nil
	}
	config := governor.RawConfig("listen.governor")
	if config != nil && !config.Enable {
		return nil
	}
	app.RegisterServers(config.Build())
	return nil
}

// isDisable
//  @Description 判断指定开关是否是关闭状态
//  @Receiver app App类型
//  @Param d 要判断的开关
//  @Return bool 开关当前状态
func (app *App) isDisable(d Disable) bool {
	b, ok := app.disableMap[d]
	if !ok {
		return false
	}
	return b
}

func (app *App) loadOrStoreTaskManager(ctx context.Context) *taskmanager.Manage {
	if app.taskManager != nil {
		return app.taskManager
	}
	manager := taskmanager.RawConfig(ctx, constant.TaskConfigKey).Build(ctx)
	app.taskManager = manager
	return app.taskManager
}
