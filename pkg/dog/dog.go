package dog

import (
	"time"
)

// A Dog is a thing that barks.
//
// A Dog is a schedule for an Idea or collection of Ideas to be "barked" to an end user.
type Dog struct {
	ID           string    `json:"id"`
	CreationTime time.Time `json:"creationTime"`
	IdeaID       string    `json:"ideaId,omitempty"`
	ScheduleType string    `json:"scheduleType"`
	Schedule     Schedule  `json:"schedule"`
	NextTaskName string    `firestore:"nextTaskName"`
	NextTaskTime time.Time `firestore:"nextTaskTime"`
}
