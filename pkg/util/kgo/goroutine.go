package kgo

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"runtime"
	"sync"
	"time"

	"github.com/codegangsta/inject"
)

// Serial
// 	@Description 生成串行执行迭代器
//	@Param fns 多个函数
// 	@Return func()  迭代器
func Serial(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

// Parallel 生成并发执行迭代器
// 	@Description: 生成并发执行迭代器
//	@Param fns 多个函数
// 	@return func() 迭代器
func Parallel(fns ...func()) func() {
	var wg sync.WaitGroup
	return func() {
		wg.Add(len(fns))
		for _, fn := range fns {
			go try2(fn, wg.Done)
		}
		wg.Wait()
	}
}

// RestrictParallel
// 	@Description: 生成并发,最大并发量restrict 迭代器
//	@Param restrict
//	@Param fns
// 	@return func() 生成迭代器
func RestrictParallel(restrict int, fns ...func()) func() {
	var channel = make(chan struct{}, restrict)
	return func() {
		var wg sync.WaitGroup
		for _, fn := range fns {
			wg.Add(1)
			go func(fn func()) {
				defer wg.Done()
				channel <- struct{}{}
				try2(fn, nil)
				<-channel
			}(fn)
		}
		wg.Wait()
		close(channel)
	}
}

// GoDirect
// 	@Description 函数执行
//	@Param fn 函数
//	@Param args 函数实参
func GoDirect(fn interface{}, args ...interface{}) {
	var inj = inject.New()
	for _, arg := range args {
		inj.Map(arg)
	}

	_, file, line, _ := runtime.Caller(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_logger.Error("recover", klog.Any("err", err), klog.String("line", fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		// 忽略返回值, goroutine执行的返回值通常都会忽略掉
		_, err := inj.Invoke(fn)
		if err != nil {
			_logger.Error("inject", klog.Any("err", err), klog.String("line", fmt.Sprintf("%s:%d", file, line)))
			return
		}
	}()
}

// Go
// 	@Description 执行go
//	@Param fn 需要运行的go
func Go(fn func()) {
	go try2(fn, nil)
}

// DelayGo
// 	@Description 延迟执行
//	@Param delay 延迟时间，time.Duration
//	@Param fn 执行函数
func DelayGo(delay time.Duration, fn func()) {
	_, file, line, _ := runtime.Caller(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_logger.Error("recover", klog.Any("err", err), klog.String("line", fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		time.Sleep(delay)
		fn()
	}()
}

// SafeGo
// 	@Description 函数
//	@Param fn 执行函数体
//	@Param rec 错误函数
func SafeGo(fn func(), rec func(error)) {
	go func() {
		err := try2(fn, nil)
		if err != nil {
			rec(err)
		}
	}()
}
