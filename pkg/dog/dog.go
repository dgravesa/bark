package dog

import (
	"encoding/json"
	"time"
)

// A Dog is a thing that barks.
//
// A Dog is a schedule for an Idea or collection of Ideas to be "barked" to an end user.
type Dog struct {
	ID           string          `json:"id" firestore:"id"`
	CreationTime time.Time       `json:"creationTime" firestore:"creationTime"`
	IdeaID       string          `json:"ideaId,omitempty" firestore:"ideaId,omitempty"`
	ScheduleType string          `json:"scheduleType" firestore:"scheduleType"`
	ScheduleRaw  json.RawMessage `json:"schedule" firestore:"schedule"`
	NextTaskName string          `json:"-" firestore:"nextTaskName"`
	NextTaskTime time.Time       `json:"-" firestore:"nextTaskTime"`

	schedule Schedule
}

// Schedule returns a Schedule based on the raw JSON in the Dog struct.
func (d *Dog) Schedule() (Schedule, error) {
	var err error = nil
	if d.schedule == nil {
		d.schedule, err = ParseSchedule("cron", d.ScheduleRaw)
	}
	return d.schedule, err
}
