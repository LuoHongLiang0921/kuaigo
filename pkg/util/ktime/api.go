// @Description time

package ktime

import "time"

// GetTimestampInMilli 获取当前时间毫秒数
// 	@Description
// 	@Return int64
func GetTimestampInMilli() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}

// Elapse 获取 函数执行耗时时间，单位 nano
// 	@Description 获取 函数执行耗时时间，单位 nano
//	@Param f
// 	@Return int64 耗时时间,单位nano 秒
func Elapse(f func()) int64 {
	now := time.Now().UnixNano()
	f()
	return time.Now().UnixNano() - now
}

// IsLeapYear ...
func IsLeapYear(year int) bool {
	if year%100 == 0 {
		return year%400 == 0
	}

	return year%4 == 0
}
