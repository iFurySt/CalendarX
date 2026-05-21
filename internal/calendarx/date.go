package calendarx

import (
	"fmt"
	"time"
)

const DateLayout = "2006-01-02"

func TodayUTC() string {
	return time.Now().UTC().Format(DateLayout)
}

func ParseDate(value string) (time.Time, error) {
	return time.Parse(DateLayout, value)
}

func MustShiftDate(value string, days int) (string, error) {
	parsed, err := ParseDate(value)
	if err != nil {
		return "", err
	}
	return parsed.AddDate(0, 0, days).Format(DateLayout), nil
}

func DateRange(anchor string, before int, after int) ([]string, error) {
	if before < 0 || after < 0 {
		return nil, fmt.Errorf("before and after must be non-negative")
	}
	start, err := ParseDate(anchor)
	if err != nil {
		return nil, err
	}
	start = start.AddDate(0, 0, -before)
	total := before + after + 1
	dates := make([]string, 0, total)
	for i := 0; i < total; i++ {
		dates = append(dates, start.AddDate(0, 0, i).Format(DateLayout))
	}
	return dates, nil
}

func CompactDate(value string) string {
	parsed, err := ParseDate(value)
	if err != nil {
		return value
	}
	return parsed.Format("20060102")
}

func NextDateCompact(value string) string {
	parsed, err := ParseDate(value)
	if err != nil {
		return value
	}
	return parsed.AddDate(0, 0, 1).Format("20060102")
}
