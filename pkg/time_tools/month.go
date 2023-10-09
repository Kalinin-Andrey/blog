package time_tools

import (
	"time"
)

func GetMonthXMonthsAgo(referenceTime time.Time, XMonthsAgo uint) *time.Time {
	t := referenceTime.Add(-time.Duration(XMonthsAgo*31*24) * time.Hour)
	ts := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return &ts
}

func GetMonthXMonthsAfter(referenceTime time.Time, XMonthsAfter uint) *time.Time {
	t := referenceTime.Add(time.Duration(XMonthsAfter*31*24) * time.Hour)
	ts := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return &ts
}
