// @Description 时间格式解析封装

package ktime

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// 常规格式化符号
var formats = map[byte]string{
	'd': "02",                        // Day:    月的哪一天，带前导零的两位数字。例：01到31
	'D': "Mon",                       // Day:    星期缩写
	'j': "2",                         // Day:    不带前导零的月份中的某一天。例：1到31
	'l': "Monday",                    // Day:    完整英文格式的周
	'F': "January",                   // Month:  完整英文格式的月
	'm': "01",                        // Month:  带前导零的月。例：01到12。
	'M': "Jan",                       // Month:  月份简写.
	'n': "1",                         // Month:  不带前导零月份
	'Y': "2006",                      // Year:   年份完整.
	'y': "06",                        // Year:   年份缩写.
	'a': "pm",                        // Time:   小写美式上午下午. Eg: am or pm.
	'A': "PM",                        // Time:   大写美式上午下午. Eg: AM or PM.
	'g': "3",                         // Time:   12小时制 不带前导零.
	'h': "03",                        // Time:   12小时制 带前导零.
	'H': "15",                        // Time:   24小时 带前导零
	'i': "04",                        // Time:   分钟 带前导零.
	's': "05",                        // Time:   分钟 不带前导零.
	'O': "-0700",                     // Zone:   与格林威治时间的时差（小时）.
	'P': "-07:00",                    // Zone:   与格林威治时间（GMT）的时差 分隔小时和分钟
	'T': "MST",                       // Zone:   Timezone abbreviation. Eg: UTC, EST, MDT ...
	'c': "2006-01-02T15:04:05-07:00", // Format: ISO 8601
	'r': "Mon, 02 Jan 06 15:04 MST",  // Format: RFC 2822 formatted date
}

// Parse
//  @Description  将标准格式时间字符串解析成 time 实例
//  @Receiver t
//  @Param value
//  @Return *Time
func (t *Time) Parse(value string) *Time {
	//if t.Error != nil {
	//	return t
	//}

	layout := DateTimeFormat

	if value == "" || value == "0" || value == "0000-00-00 00:00:00" || value == "0000-00-00" || value == "00:00:00" {
		return t
	}

	if len(value) == 10 && strings.Count(value, "-") == 2 {
		layout = DateFormat
	}

	if strings.Index(value, "T") == 10 {
		layout = RFC3339Format
	}

	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		switch len(value) {
		case 8:
			layout = ShortDateFormat
		case 14:
			layout = ShortDateTimeFormat
		}
	}

	return t.ParseByLayout(value, layout)
}

// Parse
//  @Description  将标准格式时间字符串解析成 time 实例(默认时区)
//  @Param value
//  @Return *Time
func Parse(value string) *Time {
	return SetTimezone(Local).Parse(value)
}

// ParseByLayout
//  @Description  将布局时间字符串解析成 time 实例
//  @Receiver t
//  @Param value
//  @Param layout
//  @Return *Time
func (t *Time) ParseByLayout(value string, layout string) *Time {
	loc, _ := time.LoadLocation(Local)
	tt, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		err = errors.New("the value \"" + value + "\" can't parse string as time")
	}
	t.Time = tt
	return t
}

// ParseByLayout 将布局时间字符串解析成 time 实例(默认时区)
func ParseByLayout(value string, layout string) *Time {
	return SetTimezone(Local).ParseByLayout(value, layout)
}
