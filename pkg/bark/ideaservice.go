package bark

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxIdeaBytes = 20000

// IdeaService contains handlers for the idea service endpoints.
type IdeaService struct {
	IdeaStore IdeaStore
	Logger    interface {
		Printf(format string, v ...interface{})
	}
}

func (service *IdeaService) logf(r *http.Request, format string, v ...interface{}) {
	if service.Logger != nil {
		requestID := middleware.GetReqID(r.Context())
		if requestID == "" {
			requestID = "???"
		}
		fullFormat := fmt.Sprintf("request_id=%s %s", requestID, format)
		service.Logger.Printf(fullFormat, v...)
	}
}

// IdeaStore is a data interface to all ideas.
type IdeaStore interface {
	Get(ctx context.Context, ID string) (*Idea, error)
	Put(ctx context.Context, idea *Idea) error
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

	service.logf(r, "request=get_idea idea_id=%s", id)

	idea, err := service.IdeaStore.Get(r.Context(), id)
	switch status.Code(err) {
	case codes.OK:
		service.logf(r, "idea_id=%s result=get_success", id)
		respondSuccess(w, http.StatusOK, idea)
	case codes.NotFound:
		service.logf(r, "idea_id=%s result=not_found_error", id)
		respondError(w, http.StatusNotFound, fmt.Sprint("idea not found with ID: ", id))
	default:
		service.logf(r, `idea_id=%s result=internal_error error_text="%s"`, id, err)
		respondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}

type newIdeaRequest struct {
	Text string `json:"text"`
}

// PostIdea is the handler for creating a new idea.
func (service *IdeaService) PostIdea(w http.ResponseWriter, r *http.Request) {
	var requestBody newIdeaRequest

	r.Body = http.MaxBytesReader(w, r.Body, maxIdeaBytes)
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		service.logf(r, `result=decode_error error_text="%s"`, err)
		respondError(w, http.StatusBadRequest, "could not read idea text")
		return
	} else if len(requestBody.Text) == 0 {
		service.logf(r, `result=empty_idea_error`)
		respondError(w, http.StatusBadRequest, "empty idea is not allowed")
		return
	}

	idea := &Idea{
		ID:           uuid.NewString(),
		Text:         requestBody.Text,
		CreationTime: time.Now(),
	}

	service.logf(r, `request=post_idea idea_id=%s idea_text="%s"`,
		idea.ID, html.EscapeString(idea.Text))

	err = service.IdeaStore.Put(r.Context(), idea)
	if err != nil {
		service.logf(r, `idea_id=%s result=internal_error error_text="%s"`, idea.ID, err)
		respondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	service.logf(r, "idea_id=%s result=post_success", idea.ID)
	respondSuccess(w, http.StatusCreated, idea)
}
