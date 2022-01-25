// @Description time 修饰

package ktime

import (
	"errors"
	"time"
)

// 时区常量
const (
	Local = "Local"
	CET   = "CET"
	EET   = "EET"
	EST   = "EST"
	GMT   = "GMT"
	UTC   = "UTC"
	UCT   = "UCT"
	MST   = "MST"

	Cuba      = "Cuba"
	Egypt     = "Egypt"
	Eire      = "Eire"
	Greenwich = "Greenwich"
	Iceland   = "Iceland"
	Iran      = "Iran"
	Israel    = "Israel"
	Jamaica   = "Jamaica"
	Japan     = "Japan"
	Libya     = "Libya"
	Poland    = "Poland"
	Portugal  = "Portugal"
	PRC       = "PRC"
	Singapore = "Singapore"
	Turkey    = "Turkey"
	Zulu      = "Zulu"

	Shanghai   = "Asia/Shanghai"
	Chongqing  = "Asia/Chongqing"
	HongKong   = "Asia/Hong_Kong"
	Macao      = "Asia/Macao"
	Taipei     = "Asia/Taipei"
	Tokyo      = "Asia/Tokyo"
	London     = "Europe/London"
	NewYork    = "America/New_York"
	LosAngeles = "America/Los_Angeles"
)

// 数字常量
const (
	YearsPerMillennium         = 1000    // 每千年1000年
	YearsPerCentury            = 100     // 每世纪100年
	YearsPerDecade             = 10      // 每十年10年
	QuartersPerYear            = 4       // 每年4季度
	MonthsPerYear              = 12      // 每年12月
	MonthsPerQuarter           = 3       // 每季度3月
	WeeksPerNormalYear         = 52      // 每常规年52周
	weeksPerLongYear           = 53      // 每长年53周
	WeeksPerMonth              = 4       // 每月4周
	DaysPerLeapYear            = 366     // 每闰年366天
	DaysPerNormalYear          = 365     // 每常规年365天
	DaysPerWeek                = 7       // 每周7天
	HoursPerWeek               = 168     // 每周168小时
	HoursPerDay                = 24      // 每天24小时
	MinutesPerDay              = 1440    // 每天1440分钟
	MinutesPerHour             = 60      // 每小时60分钟
	SecondsPerWeek             = 604800  // 每周604800秒
	SecondsPerDay              = 86400   // 每天86400秒
	SecondsPerHour             = 3600    // 每小时3600秒
	SecondsPerMinute           = 60      // 每分钟60秒
	MillisecondsPerSecond      = 1000    // 每秒1000毫秒
	MicrosecondsPerMillisecond = 1000    // 每毫秒1000微秒
	MicrosecondsPerSecond      = 1000000 // 每秒1000000微秒
)

// 时间格式化常量
const (
	AnsicFormat         = time.ANSIC
	UnixDateFormat      = time.UnixDate
	RubyDateFormat      = time.RubyDate
	RFC822Format        = time.RFC822
	RFC822ZFormat       = time.RFC822Z
	RFC850Format        = time.RFC850
	RFC1123Format       = time.RFC1123
	RFC1123ZFormat      = time.RFC1123Z
	RssFormat           = time.RFC1123Z
	RFC2822Format       = time.RFC1123Z
	RFC3339Format       = time.RFC3339
	KitchenFormat       = time.Kitchen
	CookieFormat        = "Monday, 02-Jan-2006 15:04:05 MST"
	RFC1036Format       = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC7231Format       = "Mon, 02 Jan 2006 15:04:05 GMT"
	DayDateTimeFormat   = "Mon, Aug 2, 2006 3:04 PM"
	DateTimeFormat      = "2006-01-02 15:04:05"
	DateFormat          = "2006-01-02"
	TimeFormat          = "15:04:05"
	UnixTimeUnitOffset  = uint64(time.Millisecond / time.Nanosecond)
	ShortDateTimeFormat = "20060102150405"
	ShortDateFormat     = "20060102"
	ShortTimeFormat     = "150405"
)

// Time time
type Time struct {
	time.Time
}

// Now returns current time
func Now() *Time {
	return &Time{
		Time: time.Now(),
	}
}

// Unix returns time converted from timestamp
func Unix(sec, nsec int64) *Time {
	return &Time{
		Time: time.Unix(sec, nsec),
	}
}

// Today
//  @Description  今天
//  @Return *Time
func Today() *Time {
	return Now().BeginOfDay()
}

// BeginOfYear
//  @Description  年开始时间
//  @Receiver t
//  @Return *Time
func (t *Time) BeginOfYear() *Time {
	y, _, _ := t.Date()
	return &Time{time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())}
}

// EndOfYear
//  @Description  年截至时间
//  @Receiver t
//  @Return *Time
func (t *Time) EndOfYear() *Time {
	return &Time{t.BeginOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)}
}

// BeginOfMonth
//  @Description  月起止时间
//  @Receiver t
//  @Return *Time
func (t *Time) BeginOfMonth() *Time {
	y, m, _ := t.Date()
	return &Time{time.Date(y, m, 1, 0, 0, 0, 0, t.Location())}
}

// EndOfMonth
//  @Description  月截至时间
//  @Receiver t
//  @Return *Time
func (t *Time) EndOfMonth() *Time {
	return &Time{t.BeginOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)}
}

// BeginOfWeek
//  @Description  周起止时间（第一天为周日）
//  @Receiver t
//  @Return *Time
func (t *Time) BeginOfWeek() *Time {
	y, m, d := t.AddDate(0, 0, 0-int(t.BeginOfDay().Weekday())).Date()
	return &Time{time.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// EndOfWeek
//  @Description  周截至时间（第一天为周日）
//  @Receiver t
//  @Return *Time
func (t *Time) EndOfWeek() *Time {
	y, m, d := t.BeginOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond).Date()
	return &Time{time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfDay
// 	@Description 开始时间
// 	@Receiver t
// 	@Return *Time
func (t *Time) BeginOfDay() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// EndOfDay
// 	@Description 结束时间
// 	@Receiver t
// 	@Return *Time
// EndOfDay returns last point of time's day
func (t *Time) EndOfDay() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfHour returns zero point of time's day
func (t *Time) BeginOfHour() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())}
}

// EndOfHour returns last point of time's day
func (t *Time) EndOfHour() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfMinute returns zero point of time's day
func (t *Time) BeginOfMinute() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, t.Location())}
}

// EndOfMinute returns last point of time's day
func (t *Time) EndOfMinute() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), t.Minute(), 59, int(time.Second-time.Nanosecond), t.Location())}
}

// ToTimestamp ToTimestampWithSecond的简称
func (t *Time) ToTimestamp() int64 {
	return t.ToTimestampWithSecond()
}

// ToTimestampWithSecond 输出秒级时间戳
func (t *Time) ToTimestampWithSecond() int64 {
	return t.Time.Unix()
}

// ToTimestampWithMillisecond 输出毫秒级时间戳
func (t *Time) ToTimestampWithMillisecond() int64 {
	return t.Time.UnixNano() / int64(time.Millisecond)
}

// ToTimestampWithMicrosecond 输出微秒级时间戳
func (t *Time) ToTimestampWithMicrosecond() int64 {
	return t.Time.UnixNano() / int64(time.Microsecond)
}

// ToTimestampWithNanosecond 输出纳秒级时间戳
func (t *Time) ToTimestampWithNanosecond() int64 {
	return t.Time.UnixNano()
}

// ToString 输出"2006-01-02 15:04:05.999999999 -0700 MST"格式字符串
func (t *Time) ToString() string {
	return t.Time.String()
}

var TS GOTimeFormat = "2006-01-02 15:04:05"

type GOTimeFormat string

func (ts GOTimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}

//const (
//	DateFormat         = "2006-01-02"
//	UnixTimeUnitOffset = uint64(time.Millisecond / time.Nanosecond)
//)

// FormatTimeMillis formats Unix timestamp (ms) to time string.
func FormatTimeMillis(tsMillis uint64) string {
	return time.Unix(0, int64(tsMillis*UnixTimeUnitOffset)).Format(string(TS))
}

// FormatDate formats Unix timestamp (ms) to date string
func FormatDate(tsMillis uint64) string {
	return time.Unix(0, int64(tsMillis*UnixTimeUnitOffset)).Format(DateFormat)
}

// CurrentTimeMillis Returns the current Unix timestamp in milliseconds.
func CurrentTimeMillis() uint64 {
	// Read from cache first.
	tickerNow := CurrentTimeMillsWithTicker()
	if tickerNow > uint64(0) {
		return tickerNow
	}
	return uint64(time.Now().UnixNano()) / UnixTimeUnitOffset
}

// Returns the current Unix timestamp in nanoseconds.
func CurrentTimeNano() uint64 {
	return uint64(time.Now().UnixNano())
}

// TODO Timezone 设置时区
func SetTimezone(name string) *Time {
	return &Time{}
}

// getLocationByTimezone 通过时区获取Location实例
func getLocationByTimezone(timezone string) (*time.Location, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		err = errors.New("invalid timezone \"" + timezone + "\", please see the $GOROOT/lib/time/zoneinfo.zip file for all valid timezone")
	}
	return loc, err
}
