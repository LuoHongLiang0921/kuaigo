// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package ktime

import (
	"sync/atomic"
	"time"
)

var nowInMs = uint64(0)

// StartTimeTicker starts a background task that caches current timestamp per millisecond,
// which may provide better performance in high-concurrency scenarios.
func StartTimeTicker() {
	atomic.StoreUint64(&nowInMs, uint64(time.Now().UnixNano())/UnixTimeUnitOffset)
	go func() {
		for {
			now := uint64(time.Now().UnixNano()) / UnixTimeUnitOffset
			atomic.StoreUint64(&nowInMs, now)
			time.Sleep(time.Millisecond)
		}
	}()
}

func CurrentTimeMillsWithTicker() uint64 {
	return atomic.LoadUint64(&nowInMs)
}
