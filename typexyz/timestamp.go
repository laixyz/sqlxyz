// Package typexyz
// 方便将数据库中整型字段值直接转换成Timestamp
// 开发者: 无双 (luciferlai@qq.com)

package typexyz

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// Timestamp 用于更快捷的使用以int64类型存储的unix时间戳
type Timestamp struct {
	Int64     int64
	TimeValue time.Time
}

//时间显示格式定义
var (
	FormatTime       = "15:04:05"
	FormatDate       = "2006-01-02"
	FormatDateZH     = "2006年01月02日"
	FormatDateTime   = "2006-01-02 15:04:05"
	FormatDateTimeZH = "2006年01月02日 15:04:05"
	FormatDateTimeT  = "2006-01-02T15:04:05"
	FormatUnknown    = "未知时间"
	FormatJust       = "刚刚"
	FormatNow        = "现在"
	FormatZHMinute   = "%d分钟前"
	FormatZHHour     = "%d小时前"
	FormatZHTMonth   = "本月02日 15:04:05"
	FormatZHMonth    = "01月02日 15:04:05"
	FormatZHLastYear = "去年01月02日 15:04:05"
)

// LayoutMode 格式模式
var LayoutMode int = 3

// String 输出字符串 默认格式
func (T Timestamp) String() string {
	if LayoutMode == 0 {
		return T.TimeValue.String()
	} else if LayoutMode == 1 {
		return T.TimeValue.Format(FormatDateTime)
	} else if LayoutMode == 2 {
		return T.TimeValue.Format(FormatDateTimeZH)
	} else if LayoutMode == 3 {
		return T.ToZhString()
	}
	return T.TimeValue.Format(FormatDateTimeT)
}

// Format 输出指定格式的时间字符串
func (T Timestamp) Format(layout string) string {
	return T.TimeValue.Format(layout)
}

// Value 用于mysql int类型字段存取unix_timestamp值
func (T Timestamp) Value() (driver.Value, error) {
	return T.Int64, nil
}

// Scan 用于mysql int类型字段存取unix_timestamp值
func (T *Timestamp) Scan(src interface{}) error {
	if src == nil {
		T.TimeValue = time.Unix(0, 0)
		T.Int64 = T.TimeValue.Unix()
		return nil
	}
	var nullInt64 sql.NullInt64
	err := nullInt64.Scan(src)
	if err != nil {
		return err
	}
	T.TimeValue = time.Unix(nullInt64.Int64, 0)
	T.Int64 = T.TimeValue.Unix()
	return nil
}

// IsExpired 是否过期，可指定时间，不指定时间则与当前时间作比较
func (T Timestamp) IsExpired(v ...time.Time) bool {
	if len(v) == 0 {
		return T.Int64 < time.Now().Unix()
	}
	return T.Int64 < v[0].Unix()
}

// ToTime 转换成Time类型
func (T Timestamp) ToTime() time.Time {
	return T.TimeValue
}

// Unix 转换成unix_stamp int64类型值
func (T Timestamp) Unix() int64 {
	return T.Int64
}

// Year 获取年份
func (T Timestamp) Year() int {
	return T.TimeValue.Year()
}

// Month 获取月份
func (T Timestamp) Month() time.Month {
	return T.TimeValue.Month()
}

// Day 获取日期day
func (T Timestamp) Day() int {
	return T.TimeValue.Day()
}

// Hour 获取日期中的Hour
func (T Timestamp) Hour() int {
	return T.TimeValue.Hour()
}

// Minute 获取日期中的 Minute
func (T Timestamp) Minute() int {
	return T.TimeValue.Minute()
}

// Second 获取日期中的 Second
func (T Timestamp) Second() int {
	return T.TimeValue.Second()
}

// ToZhString 输出具有中国特色的格式，例现在，刚刚，几分钟前，几小时前之类的
func (T Timestamp) ToZhString() string {
	if T.Int64 == 0 {
		return FormatUnknown
	}
	thisTime := time.Now()
	TTime := T.TimeValue
	sec := thisTime.Unix() - TTime.Unix()
	if sec == 0 {
		return FormatNow
	} else if sec < 60 {
		return FormatJust
	} else if sec < 3600 {
		return fmt.Sprintf(FormatZHMinute, sec/60)
	} else if sec < 86400 {
		return fmt.Sprintf(FormatZHHour, sec/3600)
	} else {
		if thisTime.Year() == TTime.Year() {
			if thisTime.Month() == TTime.Month() {
				return T.Format(FormatZHTMonth)
			}
			return T.Format(FormatZHMonth)
		} else if thisTime.Year()-1 == TTime.Year() {
			return T.Format(FormatZHLastYear)
		}
	}
	return T.Format(FormatDateTimeZH)
}

// Now 获取当前时间unix戳
func Now() Timestamp {
	var T Timestamp
	T.TimeValue = time.Now()
	T.Int64 = T.TimeValue.Unix()
	return T
}

// NewTimestamp 获取当前时间unix戳 也可指定时间戳值来创建
func NewTimestamp(v ...int64) Timestamp {
	if len(v) == 0 {
		return Now()
	}
	var T Timestamp
	T.TimeValue = time.Unix(v[0], 0)
	T.Int64 = T.TimeValue.Unix()
	return T
}

// ParseDurationToTime  类似 time.ParseDuration操作，返回参数不一样哦
func ParseDurationToTime(s string) time.Time {
	var Int64 int64
	pd, err := time.ParseDuration(s)
	if err != nil {
		Int64 = 0
	} else {
		Int64 = time.Now().Unix() + int64(pd.Seconds())
	}
	return time.Unix(Int64, 0)
}

//ParseDuration 类似 time.ParseDuration操作，返回参数不一样哦
func ParseDuration(s string) Timestamp {
	var Int64 int64
	pd, err := time.ParseDuration(s)
	if err != nil {
		Int64 = 0
	}
	Int64 = time.Now().Unix() + int64(pd.Seconds())
	return NewTimestamp(Int64)
}

// Today 获取今天的起始时间unix_timestamp值
func Today() Timestamp {
	t := time.Now()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return NewTimestamp(today.Unix())
}

// TSWK 获取本周的起始时间unix_timestamp值
func TSWK() Timestamp {
	t := time.Now()
	year, month, day := t.Date()
	var w int = int(t.Weekday())
	if w == 0 {
		w = 6
	} else {
		w = w - 1
	}
	t = time.Date(year, month, day-w, 0, 0, 0, 0, t.Location())
	return NewTimestamp(t.Unix())
}

// LastDay 获取昨天起始与结束的unix_timestamp值
func LastDay() (Timestamp, Timestamp) {
	t := time.Now()
	year, month, day := t.Date()
	begin := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	end := time.Date(year, month, day, 0, 0, -1, 0, t.Location())
	return NewTimestamp(begin.Unix()), NewTimestamp(end.Unix())
}

// LastWeek 获取上一周起始与结束的unix_timestamp值
func LastWeek() (Timestamp, Timestamp) {
	t := time.Now()
	year, month, day := t.Date()
	var w int = int(t.Weekday())
	if w == 0 {
		w = 6
	} else {
		w = w - 1
	}
	begin := time.Date(year, month, day-w-7, 0, 0, 0, 0, t.Location())
	end := time.Date(year, month, day-w, 0, 0, -1, 0, t.Location())
	return NewTimestamp(begin.Unix()), NewTimestamp(end.Unix())
}

// Date 获取指定参数的时间unix_timestamp值
func Date(year, month, day, hour, min, sec, nsec int) Timestamp {
	t := time.Now()
	t = time.Date(year, time.Month(month), day, hour, min, sec, nsec, t.Location())
	return NewTimestamp(t.Unix())
}

// ThisMonth 获取本月的起始时间unix_timestamp值
func ThisMonth() Timestamp {
	t := time.Now()
	year, month, _ := t.Date()
	today := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	return NewTimestamp(today.Unix())
}

// ThisYear 获取本月的起始时间unix_timestamp值
func ThisYear() Timestamp {
	t := time.Now()
	year, _, _ := t.Date()
	today := time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
	return NewTimestamp(today.Unix())
}
