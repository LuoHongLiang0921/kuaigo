package background

import (
	"context"
	"errors"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ktask"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"

	"github.com/pborman/uuid"

)

var (
	// MQNil mq 为空
	MQNil = errors.New("mq is nil")
)

type Handler interface {
	ktask.Handler
	GetMQ() mq.MessageQueer
}

// Background 后台任务类
type Background struct {
	config  *Config
	MQ      mq.MessageQueer
	ctx     context.Context
	logger  *klog.Logger
	handler ktask.Handler
}

// WithMQ
// 	@Description  设置消息队列
// 	@Receiver b Background
//	@Param q 消息队列
// 	@Return *Background
func (b *Background) WithMQ(q mq.MessageQueer) *Background {
	b.MQ = q
	return b
}

// WithLogger
// 	@Description 设置日志
// 	@Receiver b Background
//	@Param lg 日志
// 	@Return *Background
func (b *Background) WithLogger(lg *klog.Logger) *Background {
	b.logger = lg
	return b
}

// WithContext
// 	@Description 设置上下文
// 	@Receiver b Background
//	@Param ctx 上下文
// 	@Return *Background 设置上下文后后台任务
func (b *Background) WithContext(ctx context.Context) *Background {
	b.ctx = ctx
	return b
}

func (b *Background) Name() string {
	return b.config.Name
}

// Run
// 	@Description 启动后台任务
// 	@Receiver b Background
// 	@Return error 错误
func (b *Background) Run() error {
	if b.MQ == nil {
		return MQNil
	}
	return b.MQ.Consume(b.ctx)
}

// Stop
// 	@Description 停止后台任务
// 	@Receiver b Background
// 	@Return error 错误
func (b *Background) Stop() error {
	if b.MQ == nil {
		return MQNil
	}
	return b.MQ.Stop()
}

func (b *Background) GracefulStop() error {
	if b.MQ == nil {
		return MQNil
	}
	// 检测是否还有消息
	// 有等待，加超时
	return b.MQ.Stop()
}

// TaskType
// 	@Description 任务类型
// 	@Receiver b Background
// 	@Return string 任务类型
func (b *Background) TaskType() string {
	return constant.TaskTypeBackground
}

// RegisterHandler
// 	@Description  注册业务处理逻辑函数
// 	@Receiver b  Background
//	@Param ctx 上下文
//	@Param handler 业务逻辑处理函数
// 	@Return error 错误
func (b *Background) RegisterHandler(ctx context.Context, handler ktask.Handler) error {
	b.handler = handler
	_ = b.MQ.RegisterHandler(ctx, func(mqCtx context.Context, msg *mq.Message) error {
		if b.handler == nil {
			return fmt.Errorf("task name %v not register handler ", b.Name())
		}
		if com, ok := klog.FromContext(mqCtx); ok {
			if com.TraceId == "" {
				com.TraceId = uuid.New()
				mqCtx = klog.WithCommonLog(mqCtx, com)
			}
		}

		err := b.handler.BeforeTaskExec(mqCtx)
		if err != nil {
			return err
		}
		err = b.handler.Exec(mqCtx, msg)
		if err != nil {
			return err
		}
		_ = b.handler.AfterTaskExec(mqCtx)
		return nil
	})
	return nil
}
