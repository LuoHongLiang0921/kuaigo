// @Description 时间格式化封装

package ktime

import (
	"bytes"
	"strconv"
	"time"
)

// Format
//  @Description  ToFormatString的简称
//  @Receiver t
//  @Param format
//  @Return string
func (t *Time) Format(format string) string {
	return t.ToFormatString(format)
}

// ToFormatString
//  @Description  task_show
//  @Receiver t
//  @Param format
//  @Return string
func (t *Time) ToFormatString(format string) string {
	if t.IsZero() {
		return ""
	}
	runes := []rune(format)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(runes); i++ {
		if layout, ok := formats[byte(runes[i])]; ok {
			buffer.WriteString(t.Time.Format(layout))
		} else {
			switch runes[i] {
			case '\\': // 原样输出，不解析
				buffer.WriteRune(runes[i+1])
				i += 2
				continue
			case 'W': // ISO-8601 格式数字表示的年份中的第几周，取值范围 1-52
				buffer.WriteString(strconv.Itoa(t.WeekOfYear()))
			case 'N': // ISO-8601 格式数字表示的星期中的第几天，取值范围 1-7
				buffer.WriteString(strconv.Itoa(t.DayOfWeek()))
			case 'S': // 月份中第几天的英文缩写后缀，如st, nd, rd, th
				suffix := "th"
				switch t.Day() {
				case 1, 21, 31:
					suffix = "st"
				case 2, 22:
					suffix = "nd"
				case 3, 23:
					suffix = "rd"
				}
				buffer.WriteString(suffix)
			case 'L': // 是否为闰年，如果是闰年为 1，否则为 0
				if t.IsLeapYear() {
					buffer.WriteString("1")
				} else {
					buffer.WriteString("0")
				}
			case 'G': // 数字表示的小时，24 小时格式，没有前导零，取值范围 0-23
				buffer.WriteString(strconv.Itoa(t.Hour()))
			case 'U': // 秒级时间戳，如 1611818268
				buffer.WriteString(strconv.FormatInt(t.ToTimestamp(), 10))
			case 'u': // 数字表示的毫秒，如 999
				buffer.WriteString(strconv.Itoa(t.Millisecond()))
			case 'w': // 数字表示的星期中的第几天，取值范围 0-6
				buffer.WriteString(strconv.Itoa(t.DayOfWeek() - 1))
			case 't': // 指定的月份有几天，取值范围 28-31
				buffer.WriteString(strconv.Itoa(t.DaysInMonth()))
			case 'z': // 年份中的第几天，取值范围 0-365
				buffer.WriteString(strconv.Itoa(t.DayOfYear() - 1))
			default:
				buffer.WriteRune(runes[i])
			}
		}
	}
	return buffer.String()
}


// Tomorrow
//  @Description  明天
//  @Receiver t
//  @Return *Time
func (t *Time) Tomorrow() *Time {
	if t.Time.IsZero() {
		t.Time = time.Now().AddDate(0, 0, 1)
	} else {
		t.Time = t.Time.AddDate(0, 0, 1)
	}
	return &Time{Time:t.Time}
}

// Yesterday
//  @Description  昨天
//  @Receiver t
//  @Return *Time
func (t *Time) Yesterday() *Time {
	if t.IsZero() {
		t.Time = time.Now().AddDate(0, 0, -1)
	} else {
		t.Time = t.Time.AddDate(0, 0, -1)
	}
	return &Time{Time:t.Time}
}

// ToDateTimeString
//  @Description  输出日期时间字符串
//  @Receiver t
//  @Return string
func (t *Time) ToDateTimeString() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format(DateTimeFormat)
}

// ToDateString
//  @Description  输出日期字符串
//  @Receiver t
//  @Return string
func (t *Time) ToDateString() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format(DateFormat)
}
