// @Description 并发执行

package kgo

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// ParallelWithError
// 	@Description 生成 并行执行函数，等待所有函数执行完毕
//	@Param fns 执行函数列表
// 	@Return func() error 函数列表最后一个错误
func ParallelWithError(fns ...func() error) func() error {
	return func() error {
		eg := errgroup.Group{}
		for _, fn := range fns {
			eg.Go(fn)
		}

		return eg.Wait()
	}
}

// ParallelWithErrorChan
// 	@Description: 并行执行函数，返回所有函数执行错误
//	@Param fns 执行函数列表
// 	@return chan 错误缓冲通道
func ParallelWithErrorChan(fns ...func() error) chan error {
	total := len(fns)
	errs := make(chan error, total)

	var wg sync.WaitGroup
	wg.Add(total)

	go func(errs chan error) {
		wg.Wait()
		close(errs)
	}(errs)

	for _, fn := range fns {
		go func(fn func() error, errs chan error) {
			defer wg.Done()
			errs <- try(fn, nil)
		}(fn, errs)
	}

	return errs
}

// RestrictParallelWithErrorChan
// 	@Description: 并行执行函数,并限制运行函数个数，返回所有函数执行错误
//	@Param concurrency 限制同时运行函数的个数，个数在 1和 fns 执行函数列表长度之间
//	@Param fns 执行函数列表
// 	@Return chan 错误缓冲通道
func RestrictParallelWithErrorChan(concurrency int, fns ...func() error) chan error {
	total := len(fns)
	if concurrency <= 0 {
		concurrency = 1
	}
	if concurrency > total {
		concurrency = total
	}
	var wg sync.WaitGroup
	errs := make(chan error, total)
	jobs := make(chan func() error, concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		//consumer
		go func(jobs chan func() error, errs chan error) {
			defer wg.Done()
			for fn := range jobs {
				errs <- try(fn, nil)
			}
		}(jobs, errs)
	}
	go func(errs chan error) {
		//producer
		for _, fn := range fns {
			jobs <- fn
		}
		close(jobs)
		//wait for block errs
		wg.Wait()
		close(errs)
	}(errs)
	return errs
}
