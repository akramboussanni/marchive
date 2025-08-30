package utils

import (
	"time"
)

func IsNewUTCDay() bool {
	now := time.Now().UTC()
	return now.Hour() == 0 && now.Minute() == 0
}

func GetCurrentUTCDate() string {
	return time.Now().UTC().Format("2006-01-02")
}

func GetNextUTCMidnight() time.Time {
	now := time.Now().UTC()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return nextMidnight
}

func GetTimeUntilNextUTCMidnight() time.Duration {
	return GetNextUTCMidnight().Sub(time.Now().UTC())
}
