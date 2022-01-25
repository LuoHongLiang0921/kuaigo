package taskmanager

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask/background"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask/kcron"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask/kjob"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
)

// Manage 管理者
type Manage struct {
	// Configs 所有任务配置, key 为 taskType，value 为具体的值
	Configs map[string]*TaskConfig
	mu      sync.RWMutex
	// key wei name
	tasksByName map[string]ktask.Tasker
	tasksByType map[string][]ktask.Tasker
	tasks       []ktask.Tasker
	loadOnce    sync.Once
}

// Load
// 	@Description 载入配置
// 	@Receiver m Manage
func (m *Manage) Load(ctx context.Context) *Manage {
	m.loadOnce.Do(func() {
		for _, v := range m.Configs {
			var t ktask.Tasker
			switch v.TaskType {
			case constant.TaskTypeCron:
				t = m.doRegisterCronTask(ctx, v)
			case constant.TaskTypeOnce:
				t = m.doRegisterOnceTask(ctx, v)
			case constant.TaskTypeBackground:
				t = m.doRegisterBackgroundTask(ctx, v)
			}
			m.tasks = append(m.tasks, t)
		}
	})
	return m
}

func (m *Manage) doRegisterCronTask(ctx context.Context, c *TaskConfig) ktask.Tasker {
	config := kcron.Config{
		Name:              c.Name,
		Spec:              c.Spec,
		IsWithSeconds:     c.IsWithSeconds,
		IsImmediatelyRun:  c.IsImmediatelyRun,
		IsDistributedTask: c.IsDistributedTask,
		DelayExecType:     c.DelayExecType,
	}
	cronTask := config.WithLogger(klog.KuaigoLogger).Build(ctx)
	m.addTask(cronTask, c)
	return cronTask
}

func (m *Manage) doRegisterOnceTask(ctx context.Context, c *TaskConfig) ktask.Tasker {
	onceTask := kjob.Config{
		Name: c.Name,
	}.Build(ctx)
	m.addTask(onceTask, c)
	return onceTask
}

func (m *Manage) doRegisterBackgroundTask(ctx context.Context, c *TaskConfig) ktask.Tasker {
	backgroundTask := background.Config{
		Name: c.Name,
	}.Build().WithContext(ctx)
	m.addTask(backgroundTask, c)
	return backgroundTask
}

func (m *Manage) addTask(t ktask.Tasker, c *TaskConfig) {
	m.mu.Lock()
	m.tasksByName[c.Name] = t
	if v, ok := m.tasksByType[c.TaskType]; ok {
		v = append(v, t)
		m.tasksByType[c.TaskType] = v
	} else {
		var cronTasks []ktask.Tasker
		cronTasks = append(cronTasks, t)
		m.tasksByType[c.TaskType] = cronTasks
	}
	m.mu.Unlock()
}

// GetTaskByName
// 	@Description 根据 任务名字获取任务，配置文件中的name 字段
// 	@Receiver m Manage
//	@Param name 任务名字
// 	@Return Tasker 任务名字对应任务实例
func (m *Manage) GetTaskByName(name string) ktask.Tasker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tasksByName[name]
}

// GetTasksByType
// 	@Description 获取对应任务类型下的任务列表
// 	@Receiver m Manage
//	@Param taskType 任务类型
// 	@Return []Tasker 任务列表
func (m *Manage) GetTasksByType(taskType string) []ktask.Tasker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tasksByType[taskType]
}

// GetTasks
// 	@Description 获取所有任务列表
// 	@Receiver m Manage
// 	@Return []xtask.Tasker
func (m *Manage) GetTasks() []ktask.Tasker {
	return m.tasks
}

// RegisterHandlers
// 	@Description 注册处理接口，根据 xtask.Handler Name() 与配置文件中的name 字段匹配任务
// 	@Receiver m Manage
//	@Param ctx 上下文
//	@Param handlers 业务处理接口
func (m *Manage) RegisterHandlers(ctx context.Context, handlers []ktask.Handler) {
	for _, h := range handlers {
		taskName := h.Name()
		task := m.GetTaskByName(taskName)
		if task == nil {
			klog.KuaigoLogger.Warnf("task manager task name %s,not found task", taskName)
			continue
		}
		if task != nil {
			if task.TaskType() == constant.TaskTypeBackground {
				if backgroundTaskHandler, ok := h.(background.Handler); ok {
					if backgroundTask, ok := task.(*background.Background); ok {
						taskMQ := backgroundTaskHandler.GetMQ()
						if taskMQ != nil {
							backgroundTask.WithMQ(taskMQ).RegisterHandler(ctx, h)
						}
					}
				}
			} else {
				task.RegisterHandler(ctx, h)
			}
		}
	}
}
