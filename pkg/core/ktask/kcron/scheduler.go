package kcron

import (
	"sync/atomic"
	"time"
)

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next
// 	@Description 返回下一个激活时间
// 	@Receiver is immediatelyScheduler
//	@Param curr 当前调度时间
// 	@Return next 下一个激活时间
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}
