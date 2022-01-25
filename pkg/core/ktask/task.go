package ktask

import "context"

// Tasker 任务接口
type Tasker interface {
	Namer
	Runner
	Stoper
	GracefulStoper
	TaskTyper
	RegisterHandler(ctx context.Context, handler Handler) error
}

// Runner 执行任务
type Runner interface {
	Run() error
}

// Stoper 停止
type Stoper interface {
	Stop() error
}

// GracefulStoper 优雅停止
type GracefulStoper interface {
	GracefulStop() error
}

// Namer 任务
type Namer interface {
	Name() string
}

type TaskTyper interface {
	TaskType() string
}

// BeforeHandler 任务执行前
type BeforeHandler interface {
	// BeforeTaskExec Exec 执行逻辑
	BeforeTaskExec(ctx context.Context) error
}

// Handler 业务执行接口，业务实现这个接口
type Handler interface {
	Namer
	BeforeHandler
	// Exec 执行逻辑
	Exec(ctx context.Context, args ...interface{}) error
	AfterHandler
}

// AfterHandler 任务执行成功后
type AfterHandler interface {
	// AfterTaskExec  执行逻辑
	AfterTaskExec(ctx context.Context) error
}
