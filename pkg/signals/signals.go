// @Description

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

//Shutdown suport twice signal must exit
// Shutdown
// 	@Description  如果不等于SIGQUIT 优雅重启，两次触发直接退出应用
//  信号支持 {syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}
//	@param stop 停止后执行函数
func Shutdown(stop func(grace bool)) {
	sig := make(chan os.Signal, 2)
	signal.Notify(
		sig,
		shutdownSignals...,
	)
	go func() {
		s := <-sig
		go stop(s != syscall.SIGQUIT)
		<-sig
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}
