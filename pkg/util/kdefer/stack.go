package kdefer

import (
	"sync"
)

func NewStack() *DeferStack {
	return &DeferStack{
		fns: make([]func() error, 0),
		mu:  sync.RWMutex{},
	}
}

// DeferStack 栈
type DeferStack struct {
	fns []func() error
	mu  sync.RWMutex
}

// Push 添加一个执行函数
// 	@Description
// 	@receiver ds
//	@param fns
func (ds *DeferStack) Push(fns ...func() error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.fns = append(ds.fns, fns...)
}

// Execute
// 	@Description 执行顺序为，后进先执行
// 	@receiver ds
// todo: pop and run
func (ds *DeferStack) Execute() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i := len(ds.fns) - 1; i >= 0; i-- {
		_ = ds.fns[i]()
	}
}
