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

		var paramFields string
		params := chi.RouteContext(r.Context()).URLParams
		for i := 0; i < len(params.Keys); i++ {
			paramFields += fmt.Sprintf(" %s=%s", params.Keys[i], params.Values[i])
		}

		defaultFields := fmt.Sprintf("requestID=%s method=%s endpoint=%s%s",
			requestID, r.Method, r.URL.EscapedPath(), paramFields)

		if format == "" {
			service.Logger.Printf(defaultFields)
		} else {
			fullFormat := defaultFields + " " + format
			service.Logger.Printf(fullFormat, v...)
		}
	}
}

// IdeaStore is a data interface to all ideas.
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
		service.logf(r, "result=OK")
		respondSuccess(w, http.StatusOK, idea)
	case codes.NotFound:
		service.logf(r, "result=NotFoundError")
		respondError(w, http.StatusNotFound, fmt.Sprint("idea not found with ID: ", id))
	default:
		service.logf(r, `result=InternalError errorText="%s"`, id, err)
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
		service.logf(r, `result=DecodeError errorText="%s"`, err)
		respondError(w, http.StatusBadRequest, "could not read idea text")
		return
	} else if len(requestBody.Text) == 0 {
		service.logf(r, `result=EmptyIdeaError`)
		respondError(w, http.StatusBadRequest, "empty idea is not allowed")
		return
	}

	idea := &Idea{
		ID:           uuid.NewString(),
		Text:         requestBody.Text,
		CreationTime: time.Now(),
	}

	err = service.IdeaStore.Put(r.Context(), idea)
	if err != nil {
		service.logf(r, `ideaID=%s result=InternalError errorText="%s"`, idea.ID, err)
		respondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	service.logf(r, `ideaID=%s ideaText="%s" result=OK`, idea.ID, html.EscapeString(idea.Text))
	respondSuccess(w, http.StatusCreated, idea)
}

// DeleteIdea is the handler for deleting an idea.
func (service *IdeaService) DeleteIdea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ideaID")

	err := service.IdeaStore.Delete(r.Context(), id)
	switch status.Code(err) {
	case codes.OK:
		service.logf(r, "ideaID=%s result=OK", id)
		respondSuccess(w, http.StatusNoContent, nil)
	default:
		service.logf(r, `ideaID=%s result=InternalError errorText="%s"`, id, err)
		respondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}
