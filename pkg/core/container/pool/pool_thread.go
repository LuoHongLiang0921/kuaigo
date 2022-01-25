package pool

import "sync"

// ThreadPool :  线程池控制
type ThreadPool struct {
	queue chan struct{}
	wg    *sync.WaitGroup
}

// NewThreadPool : 新建线程池, size:goroutine个数
func NewThreadPool(size int) *ThreadPool {
	if size <= 0 {
		size = 1
	}
	return &ThreadPool{
		queue: make(chan struct{}, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add : 新加线程(routine)
func (p *ThreadPool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- struct{}{}
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done : 线程结束
func (p *ThreadPool) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait : 等待pool中所有routine结束
func (p *ThreadPool) Wait() {
	p.wg.Wait()
}
