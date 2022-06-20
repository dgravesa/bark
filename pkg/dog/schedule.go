package dog

import (
	"encoding/json"
	"errors"
	"time"
)

// Schedule defines when a Dog barks.
type Schedule interface {
	Next(time.Time) time.Time
}

// ErrScheduleTypeNotSupported is returned by ParseSchedule when the requested schedule type is
// not supported.
var ErrScheduleTypeNotSupported = errors.New("schedule type not supported")

// ParseSchedule parses raw JSON for a given schedule type.
func ParseSchedule(scheduleType string, b json.RawMessage) (Schedule, error) {
	switch scheduleType {
	case "cron":
		var schedule CronSchedule
		err := json.Unmarshal(b, &schedule)
		return schedule, err
	default:
		return nil, ErrScheduleTypeNotSupported
	}
}
