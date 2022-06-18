package bark

import (
	"context"

	"cloud.google.com/go/firestore"
)

// IdeaFirestore is a Google Cloud Firestore-based data backend for ideas.
type IdeaFirestore struct {
	FirestoreClient *firestore.Client
}

// Get returns an idea by ID.
func (store *IdeaFirestore) Get(ctx context.Context, ID string) (*Idea, error) {
	// get document from datastore
	docID := "idea/" + ID
	ideaDoc, err := store.FirestoreClient.Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}

	// convert to idea
	var idea Idea
	err = ideaDoc.DataTo(&idea)
	return &idea, err
}

// Put inserts an idea. If there is an existing idea with the same key, it will be overwritten.
func (store *IdeaFirestore) Put(ctx context.Context, idea *Idea) error {
	docID := "idea/" + idea.ID
	_, err := store.FirestoreClient.Doc(docID).Create(ctx, idea)
	return err
}

// Delete deletes an idea.
func (store *IdeaFirestore) Delete(ctx context.Context, ID string) error {
	docID := "idea/" + ID
	_, err := store.FirestoreClient.Doc(docID).Delete(ctx)
	return err
}
