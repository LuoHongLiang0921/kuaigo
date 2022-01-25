// @Description duration

package ktime

import "time"

// Duration ...
// 	@Description panic if parse duration failed,"ns", "us" (or "µs"), "ms", "s", "m", "h"
//	@Param str
// 	@Return time.Duration
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}

	return dur
}
