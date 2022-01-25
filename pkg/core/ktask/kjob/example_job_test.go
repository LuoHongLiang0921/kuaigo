package kjob_test

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
)

// go run main.go --job=demoJob --disable-once-job=false
func main() {
	ctx := context.Background()
	eng := NewEngine(ctx)
	if err := eng.Run(); err != nil {
		klog.Fatal(err.Error())
	}
	//ctx:=context.Background()
	//onceTask := NewOnceTask("click me too")
	//tabby.GetInstance().WithContext(ctx).Build().RegisterTasks(onceTask).Run()
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
	onceTask := NewOnceTask("click me too")
	e.WithContext(ctx).RegisterTasks(onceTask)
	return nil
}

type OnceTask struct {
	Content string
}

func NewOnceTask(content string) *OnceTask {
	return &OnceTask{Content: content}
}

func (o OnceTask) Name() string {
	return "demoJob"
}

func (o OnceTask) BeforeTaskExec(ctx context.Context) error {
	klog.WithContext(ctx).Infof(" BeforeTaskExec content %v", o.Content)
	return nil
}

func (o OnceTask) Exec(ctx context.Context, args ...interface{}) error {
	klog.WithContext(ctx).Infof("content %v", o.Content)
	closed, ok := args[0].(chan struct{})
	if ok {
		<-closed
		return nil
	}

	return nil
}

func (o OnceTask) AfterTaskExec(ctx context.Context) error {
	klog.WithContext(ctx).Infof("AfterTaskExec content %v", o.Content)
	return nil
}

func ExampleConfig_Build() {
	ctx := context.Background()
	eng := NewEngine(ctx)
	if err := eng.Run(); err != nil {
		klog.WithContext(ctx).Error(err.Error())
	}
}
