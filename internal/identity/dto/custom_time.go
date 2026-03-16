package dto

import "time"

// Normal time but only ms is only precise up to 6 digits
type CustomTime struct {
	Time time.Time
}

func NewCustomTime(timeIn time.Time) CustomTime {
	return CustomTime{
		Time: timeIn.UTC().Truncate(time.Microsecond),
	}
}

func NewCustomTimeNow() CustomTime {
	return NewCustomTime(time.Now())
}
