package utils

import "time"

func Now() time.Time {
	return time.Now()
}

func NowNanoSeconds() float64 {
	now := Now()
	micSeconds := float64(now.Nanosecond()) / 1000000000

	return float64(now.Unix()) + micSeconds
}
