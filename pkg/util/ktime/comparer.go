// @Description 时间比较库

package ktime

import "time"

// IsZero
//  @Description  是否是零值
//  @Receiver t
//  @Return bool
func (t *Time) IsZero() bool {
	return t.Time.IsZero()
}

// IsNow
//  @Description  是否是当前时间
//  @Receiver t
//  @Return bool
func (t *Time) IsNow() bool {
	return t.ToTimestamp() == t.ToTimestamp()
}

// IsFuture
//  @Description  是否是未来时间
//  @Receiver t
//  @Return bool
func (t *Time) IsFuture() bool {
	return t.ToTimestamp() > t.ToTimestamp()
}

// IsPast
//  @Description  是否是过去时间
//  @Receiver t
//  @Return bool
func (t *Time) IsPast() bool {
	return t.ToTimestamp() < t.ToTimestamp()
}

// IsLeapYear
//  @Description  是否是闰年
//  @Receiver t
//  @Return bool
func (t *Time) IsLeapYear() bool {
	year := t.Time.Year()
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}


// IsJanuary
//  @Description  是否是一月
//  @Receiver t
//  @Return bool
func (t *Time) IsJanuary() bool {
	return t.Time.Month() == time.January
}

// IsFebruary
//  @Description  是否是二月
//  @Receiver t
//  @Return bool
func (t *Time) IsFebruary() bool {
	return t.Time.Month() == time.February
}

// IsMarch
//  @Description  是否是三月
//  @Receiver t
//  @Return bool
func (t *Time) IsMarch() bool {
	return t.Time.Month() == time.March
}

// IsApril
//  @Description  是否是四月
//  @Receiver t
//  @Return bool
func (t *Time) IsApril() bool {
	return t.Time.Month() == time.April
}

// IsMay
//  @Description  是否是五月
//  @Receiver t
//  @Return bool
func (t *Time) IsMay() bool {
	return t.Time.Month() == time.May
}

// IsJune
//  @Description  是否是六月
//  @Receiver t
//  @Return bool
func (t *Time) IsJune() bool {
	return t.Time.Month() == time.June
}

// IsJuly
//  @Description  是否是七月
//  @Receiver t
//  @Return bool
func (t *Time) IsJuly() bool {
	return t.Time.Month() == time.July
}

// IsAugust
//  @Description  是否是八月
//  @Receiver t
//  @Return bool
func (t *Time) IsAugust() bool {
	return t.Time.Month() == time.August
}

// IsSeptember
//  @Description  是否是九月
//  @Receiver t
//  @Return bool
func (t *Time) IsSeptember() bool {
	return t.Time.Month() == time.September
}

// IsOctober
//  @Description  是否是十月
//  @Receiver t
//  @Return bool
func (t *Time) IsOctober() bool {
	return t.Time.Month() == time.October
}

// IsNovember
//  @Description  是否是十一月
//  @Receiver t
//  @Return bool
func (t *Time) IsNovember() bool {
	return t.Time.Month() == time.November
}

// IsDecember
//  @Description  是否是十二月
//  @Receiver t
//  @Return bool
func (t *Time) IsDecember() bool {
	return t.Time.Month() == time.December
}

// IsMonday
//  @Description  是否是周一
//  @Receiver t
//  @Return bool
func (t *Time) IsMonday() bool {
	return t.Time.Weekday() == time.Monday
}

// IsTuesday
//  @Description  是否是周二
//  @Receiver t
//  @Return bool
func (t *Time) IsTuesday() bool {
	return t.Time.Weekday() == time.Tuesday
}

// IsWednesday
//  @Description  是否是周三
//  @Receiver t
//  @Return bool
func (t *Time) IsWednesday() bool {
	return t.Time.Weekday() == time.Wednesday
}

// IsThursday
//  @Description  是否是周四
//  @Receiver t
//  @Return bool
func (t *Time) IsThursday() bool {
	return t.Time.Weekday() == time.Thursday
}

// IsFriday
//  @Description  是否是周五
//  @Receiver t
//  @Return bool
func (t *Time) IsFriday() bool {
	return t.Time.Weekday() == time.Friday
}

// IsSaturday
//  @Description  是否是周六
//  @Receiver t
//  @Return bool
func (t *Time) IsSaturday() bool {
	return t.Time.Weekday() == time.Saturday
}

// IsSunday
//  @Description  是否是周日
//  @Receiver t
//  @Return bool
func (t *Time) IsSunday() bool {
	return t.Time.Weekday() == time.Sunday
}

// IsWeekday
//  @Description  是否是工作日
//  @Receiver t
//  @Return bool
func (t *Time) IsWeekday() bool {
	return !t.IsSaturday() && !t.IsSunday()
}

// IsWeekend
//  @Description  是否是周末
//  @Receiver t
//  @Return bool
func (t *Time) IsWeekend() bool {
	return t.IsSaturday() || t.IsSunday()
}


// Compare
//  @Description  时间比较
//  @Receiver t
//  @Param operator
//  @Param tt
//  @Return bool
func (t *Time) Compare(operator string, tt *Time) bool {
	switch operator {
	case "=":
		return t.Eq(tt)
	case "<>":
		return !t.Eq(tt)
	case "!=":
		return !t.Eq(tt)
	case ">":
		return t.Gt(tt)
	case ">=":
		return t.Gte(tt)
	case "<":
		return t.Lt(tt)
	case "<=":
		return t.Lte(tt)
	}

	return false
}

// Gt
//  @Description  大于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Gt(tt *Time) bool {
	return t.Time.After(tt.Time)
}

// Lt
//  @Description  小于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Lt(tt *Time) bool {
	return t.Time.Before(tt.Time)
}

// Eq
//  @Description  等于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Eq(tt *Time) bool {
	return t.Time.Equal(tt.Time)
}

// Ne
//  @Description  不等于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Ne(tt *Time) bool {
	return !t.Eq(tt)
}

// Gte
//  @Description  大于等于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Gte(tt *Time) bool {
	return t.Gt(tt) || t.Eq(tt)
}

// Lte
//  @Description  小于等于
//  @Receiver t
//  @Param tt
//  @Return bool
func (t *Time) Lte(tt *Time) bool {
	return t.Lt(tt) || t.Eq(tt)
}

// Between
//  @Description  是否在两个时间之间(不包括这两个时间)
//  @Receiver t
//  @Param start
//  @Param end
//  @Return bool
func (t *Time) Between(start *Time, end *Time) bool {
	if t.Gt(start) && t.Lt(end) {
		return true
	}
	return false
}

// BetweenIncludedStartTime
//  @Description  是否在两个时间之间(包括开始时间)
//  @Receiver t
//  @Param start
//  @Param end
//  @Return bool
func (t *Time) BetweenIncludedStartTime(start *Time, end *Time) bool {
	if t.Gte(start) && t.Lt(end) {
		return true
	}
	return false
}

// BetweenIncludedEndTime
//  @Description  是否在两个时间之间(包括结束时间)
//  @Receiver t
//  @Param start
//  @Param end
//  @Return bool
func (t *Time) BetweenIncludedEndTime(start *Time, end *Time) bool {
	if t.Gt(start) && t.Lte(end) {
		return true
	}
	return false
}

// BetweenIncludedBoth
//  @Description  是否在两个时间之间(包括这两个时间)
//  @Receiver t
//  @Param start
//  @Param end
//  @Return bool
func (t *Time) BetweenIncludedBoth(start *Time, end *Time) bool {
	if t.Gte(start) && t.Lte(end) {
		return true
	}
	return false
}
