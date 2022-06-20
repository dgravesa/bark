package bark

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxIdeaBytes = 20000

// IdeaService contains handlers for the idea service endpoints.
type IdeaService struct {
	Service
	IdeaStore IdeaStore
}

// IdeaStore is an interface to a data store for ideas.
type IdeaStore interface {
	Get(ctx context.Context, ID string) (*Idea, error)
	Put(ctx context.Context, idea *Idea) error
	Delete(ctx context.Context, ID string) error
}

// RegisterRoutes registers the idea service routes to the chi router.
func (service *IdeaService) RegisterRoutes(r *chi.Mux) {
	r.Route("/ideas", func(r chi.Router) {
		r.Post("/", service.PostIdea)

		r.Route("/{ideaID}", func(r chi.Router) {
			r.Get("/", service.GetIdea)
			r.Delete("/", service.DeleteIdea)
		})
	})
}

// GetIdea is the handler for getting a single idea by ID.
func (service *IdeaService) GetIdea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ideaID")

	idea, err := service.IdeaStore.Get(r.Context(), id)
	switch status.Code(err) {
	case codes.OK:
		service.Logf(r, "result=OK")
		RespondSuccess(w, http.StatusOK, idea)
	case codes.NotFound:
		service.Logf(r, "result=NotFoundError")
		RespondError(w, http.StatusNotFound, fmt.Sprint("idea not found with ID: ", id))
	default:
		service.Logf(r, `result=InternalError errorText="%s"`, id, err)
		RespondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}

// NewIdeaRequest is the request type for a new idea.
type NewIdeaRequest struct {
	Text string `json:"text"`
}

// PostIdea is the handler for creating a new idea.
func (service *IdeaService) PostIdea(w http.ResponseWriter, r *http.Request) {
	var requestBody NewIdeaRequest

	r.Body = http.MaxBytesReader(w, r.Body, maxIdeaBytes)
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		service.Logf(r, `result=DecodeError errorText="%s"`, err)
		RespondError(w, http.StatusBadRequest, "could not read idea text")
		return
	} else if len(requestBody.Text) == 0 {
		service.Logf(r, `result=EmptyIdeaError`)
		RespondError(w, http.StatusBadRequest, "empty idea is not allowed")
		return
	}

	idea := &Idea{
		ID:           uuid.NewString(),
		Text:         requestBody.Text,
		CreationTime: time.Now(),
	}

	err = service.IdeaStore.Put(r.Context(), idea)
	if err != nil {
		service.Logf(r, `ideaID=%s result=InternalError errorText="%s"`, idea.ID, err)
		RespondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	service.Logf(r, `ideaID=%s ideaText="%s" result=OK`, idea.ID, html.EscapeString(idea.Text))
	RespondSuccess(w, http.StatusCreated, idea)
}

// DeleteIdea is the handler for deleting an idea.
func (service *IdeaService) DeleteIdea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ideaID")

	err := service.IdeaStore.Delete(r.Context(), id)
	switch status.Code(err) {
	case codes.OK:
		service.Logf(r, "ideaID=%s result=OK", id)
		RespondSuccess(w, http.StatusNoContent, nil)
	default:
		service.Logf(r, `ideaID=%s result=InternalError errorText="%s"`, id, err)
		RespondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}
