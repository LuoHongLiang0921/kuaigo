package kjob

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"sync/atomic"
)

// XJob 一次性任务
type XJob struct {
	config  *Config
	handler ktask.Handler
	ctx     context.Context

	closed  chan struct{}
	isCosed int32
}

// Name
// 	@Description
// 	@Receiver x
// 	@Return string
func (x *XJob) Name() string {
	return x.config.Name
}

// Run
// 	@Description
// 	@Receiver x
// 	@Return error
func (x *XJob) Run() error {
	if x.handler == nil {
		return fmt.Errorf("task name %v not register handler ", x.Name())
	}

	err := x.handler.BeforeTaskExec(x.ctx)
	if err != nil {
		return err
	}
	err = x.handler.Exec(x.ctx, x.closed)
	if err != nil {
		return err
	}
	return x.handler.AfterTaskExec(x.ctx)
}

// Stop
// 	@Description
// 	@Receiver x
// 	@Return error
func (x *XJob) Stop() error {
	if x.closed != nil {
		if atomic.LoadInt32(&x.isCosed) == 1 {
			return nil
		}
		close(x.closed)
		atomic.StoreInt32(&x.isCosed, 1)
	}
	return nil
}

// GracefulStop
// 	@Description
// 	@Receiver x
// 	@Return error
func (x *XJob) GracefulStop() error {
	// 检测是否还有消息
	// 有等待，加超时
	return x.Stop()
}

// TaskType
// 	@Description
// 	@Receiver x
// 	@Return string
func (x *XJob) TaskType() string {
	return constant.TaskTypeOnce
}

// RegisterHandler
// 	@Description
// 	@Receiver x
//	@Param ctx
//	@Param handler
// 	@Return error
func (x *XJob) RegisterHandler(ctx context.Context, handler ktask.Handler) error {
	x.handler = handler
	return nil
}
