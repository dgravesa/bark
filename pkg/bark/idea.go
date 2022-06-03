package bark

import "time"

// An Idea is a thing people have when they get smart or stupid.
type Idea struct {
	ID           string    `json:"id" datastore:"id"`
	Text         string    `json:"text" datastore:"text"`
	CreationTime time.Time `json:"creationTime" datastore:"creationTime"`
}
