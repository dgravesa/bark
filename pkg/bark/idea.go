package bark

import "time"

// An Idea is a thing people have when they get smart or stupid.
type Idea struct {
	ID           string    `json:"id" datastore:"id" firestore:"id"`
	Text         string    `json:"text" datastore:"text" firestore:"text"`
	CreationTime time.Time `json:"creationTime" datastore:"creationTime" firestore:"creationTime"`
}
