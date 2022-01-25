package kuaigo

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask/taskmanager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcycle"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kdefer"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"

)

type Option func(a *App)

type Disable byte

const (
	DisableParserFlag      Disable = 1
	DisableLoadConfig      Disable = 2
	DisableDefaultGovernor Disable = 3
)

const (
	// StageAfterStart before app start
	StageAfterStart uint32 = 0
	// StageAfterStop after app stop
	StageAfterStop uint32 = 1
	// StageBeforeStop before app stop
	StageBeforeStop uint32 = 2
)

type App struct {
	cycle        *kcycle.Cycle
	smu          *sync.RWMutex
	buildOnce    sync.Once
	stopOnce     sync.Once
	initFns      []func() error
	servers      []server.Server
	taskManager  *taskmanager.Manage
	logger       *klog.Logger
	configParser conf.Unmarshaller
	disableMap   map[Disable]bool
	HideBanner   bool
	hooks        map[uint32]*kdefer.DeferStack
	// ServiceName 服务名
	ServiceName string
	// app 上下文，为了兼容老版本
	ctx context.Context
	// 版本号
	ConfigVersion string
}

var appInstance *App
var instanceOnce sync.Once
