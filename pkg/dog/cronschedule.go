package dog

import (
	"encoding/json"
	"time"

	"github.com/robfig/cron"
)

// CronSchedule represents a Cron-based schedule.
type CronSchedule struct {
	spec string
	cron.Schedule
}

// Next implements the Schedule interface.
func (s CronSchedule) Next(t time.Time) time.Time {
	schedule, _ := cron.Parse(s.spec)
	return schedule.Next(t)
}

// UnmarshalJSON parses a JSON string into a CronSchedule.
func (s *CronSchedule) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.spec)
	if err != nil {
		return err
	}
	// validate spec
	s.Schedule, err = cron.ParseStandard(s.spec)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON writes the Cron spec as a JSON string.
func (s *CronSchedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.spec)
}
