package time_util

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

const TimeFormat = "2006-01-02 15:04:05"

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	// 如果时间是零值，返回空字符串或 null
	if tTime.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format(TimeFormat))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	if tTime.IsZero() {
		return nil, nil
	}
	return tTime, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalTime", v)
}
