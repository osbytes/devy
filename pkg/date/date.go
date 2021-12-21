package date

import (
	"time"
)

// EOD returns the end of the day in the provided timezone
func EOD(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, loc)
}

// SOD returns the start of the day in the provided timezone
func SOD(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

// SOY returns the start of the year in the provided timezone
func SOY(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, loc)
}

// EOY returns the end of the year in the provided timezone
func EOY(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 12, 31, 59, 59, 0, 0, loc)
}

// WithinDuration test to see if two time.Time's are within a duration of eachother
func WithinDuration(expected, actual time.Time, delta time.Duration) bool {
	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		return false
	}

	return true
}
