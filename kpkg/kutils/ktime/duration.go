// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package ktime

import "time"

// Duration ...
// panic if parse duration failed
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}

	return dur
}
