package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/bark/pkg/bark"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxIdeaBytes = 20000

// IdeaService contains handlers for the idea service endpoints.
type IdeaService struct {
	IdeaStore IdeaStore
}

// IdeaStore is a data interface to all ideas.
type IdeaStore interface {
	Get(ctx context.Context, ID string) (*bark.Idea, error)
	Put(ctx context.Context, idea *bark.Idea) error
}

// RegisterRoutes registers the idea service routes to the chi router.
func (service *IdeaService) RegisterRoutes(r *chi.Mux) {
	r.Route("/ideas", func(r chi.Router) {
		r.Post("/", service.PostIdea)

		r.Route("/{ideaID}", func(r chi.Router) {
			r.Get("/", service.GetIdea)
		})
	})
}

// GetIdea is the handler for getting a single idea by ID.
func (service *IdeaService) GetIdea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ideaID")

	idea, err := service.IdeaStore.Get(r.Context(), id)
	switch err {
	case nil:
		respondSuccess(w, http.StatusOK, idea)
	case datastore.ErrNoSuchEntity:
		respondError(w, http.StatusNotFound, fmt.Sprint("idea not found with ID: ", id))
	default:
		respondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}

type newIdeaRequest struct {
	Text string `json:"idea"`
}

// PostIdea is the handler for creating a new idea.
func (service *IdeaService) PostIdea(w http.ResponseWriter, r *http.Request) {
	var requestBody newIdeaRequest

	r.Body = http.MaxBytesReader(w, r.Body, maxIdeaBytes)

	rd := http.MaxBytesReader(w, r.Body, maxIdeaBytes)
	d := json.NewDecoder(rd)
	err := d.Decode(&requestBody)
	if err != nil {
		respondError(w, http.StatusBadRequest, "could not read idea text")
		return
	} else if len(requestBody.Text) == 0 {
		respondError(w, http.StatusBadRequest, "empty idea is not allowed")
		return
	}

	idea := &bark.Idea{
		ID:           uuid.NewString(),
		Text:         requestBody.Text,
		CreationTime: time.Now(),
	}

	err = service.IdeaStore.Put(r.Context(), idea)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	respondSuccess(w, http.StatusCreated, idea)
}
