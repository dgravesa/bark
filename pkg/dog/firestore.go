package dog

import (
	"context"

	"cloud.google.com/go/firestore"
)

// DoggoFirestore is a Google Cloud Firestore-based data backend for doggos.
type DoggoFirestore struct {
	FirestoreClient *firestore.Client
}

// Get returns a dog by ID.
func (store *DoggoFirestore) Get(ctx context.Context, ID string) (*Dog, error) {
	// get document from datastore
	docID := "dogs/" + ID
	dogDoc, err := store.FirestoreClient.Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}

	// convert to dog
	var dog Dog
	err = dogDoc.DataTo(&dog)
	return &dog, err
}

// Put inserts an dog. If there is an existing dog with the same key, it will be overwritten.
func (store *DoggoFirestore) Put(ctx context.Context, dog *Dog) error {
	docID := "dogs/" + dog.ID
	_, err := store.FirestoreClient.Doc(docID).Create(ctx, dog)
	return err
}

// Delete deletes an dog.
func (store *DoggoFirestore) Delete(ctx context.Context, ID string) error {
	docID := "dogs/" + ID
	_, err := store.FirestoreClient.Doc(docID).Delete(ctx)
	return err
}
