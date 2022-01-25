// @Description 获取时间封装库

package ktime

// DaysInYear
//  @Description  获取本年的总天数
//  @Receiver t
//  @Return int
func (t *Time) DaysInYear() int {
	if t.IsZero() {
		return 0
	}
	days := DaysPerNormalYear
	if t.IsLeapYear() {
		days = DaysPerLeapYear
	}
	return days
}

// DaysInMonth
//  @Description  获取本月的总天数
//  @Receiver t
//  @Return int
func (t *Time) DaysInMonth() int {
	if t.IsZero() {
		return 0
	}
	return t.EndOfMonth().Time.Day()
}

// MonthOfYear
//  @Description  获取本年的第几月(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) MonthOfYear() int {
	if t.IsZero() {
		return 0
	}
	return int(t.Time.Month())
}

// DayOfYear
//  @Description  获取本年的第几天(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) DayOfYear() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.YearDay()
}

// DayOfMonth
//  @Description  获取本月的第几天(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) DayOfMonth() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Day()
}

// DayOfWeek
//  @Description  获取本周的第几天(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) DayOfWeek() int {
	if t.IsZero() {
		return 0
	}
	day := int(t.Time.Weekday())
	if day == 0 {
		return DaysPerWeek
	}
	return day
}

// WeekOfYear
//  @Description  获取本年的第几周(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) WeekOfYear() int {
	if t.IsZero() {
		return 0
	}
	_, week := t.Time.ISOWeek()
	return week
}

// WeekOfMonth
//  @Description  获取本月的第几周(从1开始)
//  @Receiver t
//  @Return int
func (t *Time) WeekOfMonth() int {
	if t.IsZero() {
		return 0
	}
	day := t.Time.Day()
	if day < DaysPerWeek {
		return 1
	}
	return day%DaysPerWeek + 1
}

// Century
//  @Description  获取当前世纪
//  @Receiver t
//  @Return int
func (t *Time) Century() int {
	if t.IsZero() {
		return 0
	}
	return t.Year()/100 + 1
}

// Year
//  @Description  获取当前年
//  @Receiver t
//  @Return int
func (t *Time) Year() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Year()
}

// Quarter
//  @Description  获取当前季度
//  @Receiver t
//  @Return int
func (t *Time) Quarter() int {
	if t.IsZero() {
		return 0
	}
	switch {
	case t.Month() >= 10:
		return 4
	case t.Month() >= 7:
		return 3
	case t.Month() >= 4:
		return 2
	case t.Month() >= 1:
		return 1
	default:
		return 0
	}
}

// Month
//  @Description  获取当前月
//  @Receiver t
//  @Return int
func (t *Time) Month() int {
	if t.IsZero() {
		return 0
	}
	return t.MonthOfYear()
}

// Week
//  @Description  获取当前周(从0开始)
//  @Receiver t
//  @Return int
func (t *Time) Week() int {
	if t.IsZero() {
		return -1
	}
	return int(t.Time.Weekday())
}

// Day
//  @Description  获取当前日
//  @Receiver t
//  @Return int
func (t *Time) Day() int {
	if t.IsZero() {
		return 0
	}
	return t.DayOfMonth()
}

// Hour
//  @Description  获取当前小时
//  @Receiver t
//  @Return int
func (t *Time) Hour() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Hour()
}

// Minute
//  @Description  获取当前分钟数
//  @Receiver t
//  @Return int
func (t *Time) Minute() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Minute()
}

// Second
//  @Description  获取当前秒数
//  @Receiver t
//  @Return int
func (t *Time) Second() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Second()
}

// Millisecond
//  @Description  获取当前毫秒数
//  @Receiver t
//  @Return int
func (t *Time) Millisecond() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Nanosecond() / 1e6
}

// Microsecond
//  @Description  获取当前微秒数
//  @Receiver t
//  @Return int
func (t *Time) Microsecond() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Nanosecond() / 1e9
}

// Nanosecond
//  @Description  获取当前纳秒数
//  @Receiver t
//  @Return int
func (t *Time) Nanosecond() int {
	if t.IsZero() {
		return 0
	}
	return t.Time.Nanosecond()
}