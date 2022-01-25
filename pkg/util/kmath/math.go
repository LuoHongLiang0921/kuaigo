// @Description math 工具

package kmath

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrConvertFail = errors.New("convert data type is failure")
)

// Abs
//  @Description  取绝对值
//  @Param number
//  @Return float64
func Abs(number float64) float64 {
	return math.Abs(number)
}

// Percent
//  @Description  返回百分比
//  @Param val 占比数字
//  @Param total 总数字
//  @Return float64
func Percent(val, total int) float64 {
	if total == 0 {
		return float64(0)
	}

	return (float64(val) / float64(total)) * 100
}

// Rand
//  @Description   取随机数
//  @Param min
//  @Param max
//  @Return int
//  @Return error
func Rand(min, max int) (int, error) {
	if min > max {
		return -1, errors.New("min: min cannot be greater than max")
	}
	if int31 := 1<<31 - 1; max > int31 {
		return -1, errors.New("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min, nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min, nil
}

// Round
//  @Description  四舍五入
//  @Param value
//  @Return float64
func Round(value float64) float64 {
	return math.Floor(value + 0.5)
}

// Floor
//  @Description  转换为浮点型
//  @Param value
//  @Return float64
func Floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil
//  @Description 向上取整
//  @Param value
//  @Return float64
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pi
//  @Description  获取圆周率
//  @Return float64
func Pi() float64 {
	return math.Pi
}

// Max
//  @Description  取区间最大值
//  @Param nums
//  @Return float64
//  @Return error
func Max(nums ...float64) (float64, error) {
	if len(nums) < 2 {
		return -1, errors.New("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max, nil
}

// Min min(),异常时返回-1，需要判断error是否为nil
func Min(nums ...float64) (float64, error) {
	if len(nums) < 2 {
		return -1, errors.New("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min, nil
}

//Fibonacci 内部调用
func fib() func() int {
	first, second := 1, 1
	return func() int {
		ret := first
		first, second = second, first+second
		return ret
	}
}

//Fibonacci 斐波那契数列
//	一般用于告警间隔时间，0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584,4181
func Fibonacci(num int) (result int) {
	f := fib()
	for i := 1; i < 1000; i++ {
		if num == i {
			return f()
		}
		f()
	}
	return
}

/*************************************************************
 * convert value to int
 *************************************************************/

// Int
//  @Description  转换为整型 带返回错误
//  @Param in
//  @Return int
//  @Return error
func Int(in interface{}) (int, error) {
	return ToInt(in)
}

// MustInt
//  @Description  转换为整型
//  @Param in
//  @Return int
func MustInt(in interface{}) int {
	val, _ := ToInt(in)
	return val
}

// ToInt
//  @Description  常见类型转换为整型
//  @Param in
//  @Return iVal
//  @Return err
func ToInt(in interface{}) (iVal int, err error) {
	switch tVal := in.(type) {
	case nil:
		iVal = 0
	case int:
		iVal = tVal
	case int8:
		iVal = int(tVal)
	case int16:
		iVal = int(tVal)
	case int32:
		iVal = int(tVal)
	case int64:
		iVal = int(tVal)
	case uint:
		iVal = int(tVal)
	case uint8:
		iVal = int(tVal)
	case uint16:
		iVal = int(tVal)
	case uint32:
		iVal = int(tVal)
	case uint64:
		iVal = int(tVal)
	case float32:
		iVal = int(tVal)
	case float64:
		iVal = int(tVal)
	case string:
		iVal, err = strconv.Atoi(strings.TrimSpace(tVal))
	default:
		err = ErrConvertFail
	}
	return
}
