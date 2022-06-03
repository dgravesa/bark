package bark

import (
	"context"

	"cloud.google.com/go/datastore"
)

// IdeaDatastore is a Google Cloud Datastore-based data backend for ideas.
type IdeaDatastore struct {
	DatastoreClient *datastore.Client
}

// Get returns an idea by ID.
func (store *IdeaDatastore) Get(ctx context.Context, ID string) (*Idea, error) {
	var idea Idea
	key := datastore.NameKey("idea", ID, nil)
	err := store.DatastoreClient.Get(ctx, key, &idea)
	return &idea, err
}

// Put inserts an idea. If there is an existing idea with the same key, it will be overwritten.
func (store *IdeaDatastore) Put(ctx context.Context, idea *Idea) error {
	key := datastore.NameKey("idea", idea.ID, nil)
	_, err := store.DatastoreClient.Put(ctx, key, idea)
	return err
}
