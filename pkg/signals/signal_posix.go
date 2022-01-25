// @Description 信号，在window 下编译

// +build !windows

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}
