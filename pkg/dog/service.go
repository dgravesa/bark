package dog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgravesa/bark/pkg/bark"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service contains handlers for the Dog scheduling service endpoints.
type Service struct {
	bark.Service
	IdeaGetter  IdeaGetter
	DogGetter   DoggoGetter
	TasksClient TasksClient
}

// IdeaGetter is an interface for getting ideas.
type IdeaGetter interface {
	Get(ctx context.Context, ID string) (*bark.Idea, error)
}

// DoggoGetter is an interface for getting doggos.
type DoggoGetter interface {
	Get(ctx context.Context, ID string) (*Dog, error)
}

// TasksClient is a client interface to tasks.
type TasksClient interface {
	Register(ctx context.Context, dog *Dog) (*Dog, error)
	Unregister(ctx context.Context, dogID string) error
}

// RegisterRoutes registers service routes to a chi mux instance.
func (service *Service) RegisterRoutes(r *chi.Mux) {
	r.Route("/dogs", func(r chi.Router) {
		r.Post("/", service.PostDog)

		r.Route("/{dogID}", func(r chi.Router) {
			r.Get("/", service.GetDog)
			r.Delete("/", service.DeleteDog)
		})
	})
}

// CreateDogRequest is the request type for creating a new Dog.
type CreateDogRequest struct {
	IdeaID       string          `json:"ideaId,omitempty"`
	ScheduleType string          `json:"scheduleType"`
	Schedule     json.RawMessage `json:"schedule"`
}

var maxCreateDogRequestSizeBytes int64 = 20000

// PostDog is a handler for creating a new Dog.
func (service *Service) PostDog(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateDogRequest

	// read request body into struct
	r.Body = http.MaxBytesReader(w, r.Body, maxCreateDogRequestSizeBytes)
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		service.Logf(r, `result=DecodeError errorText="%s"`, err)
		bark.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	ideaID := requestBody.IdeaID

	schedule, err := ParseSchedule(requestBody.ScheduleType, requestBody.Schedule)
	if err != nil {
		service.Logf(r, `result=ParseScheduleError errorText="%s"`, err)
		bark.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// verify target exists
	_, err = service.IdeaGetter.Get(r.Context(), ideaID)
	switch status.Code(err) {
	case codes.OK:
		service.Logf(r, `action=GetIdea ideaID=%s result=OK`, ideaID)
	case codes.NotFound:
		service.Logf(r, `action=GetIdea ideaID=%s result=NotFoundError`, ideaID)
		bark.RespondError(w, http.StatusNotFound, fmt.Sprint("idea not found with ID: ", ideaID))
		return
	default:
		service.Logf(r, `action=GetIdea ideaID=%s result=InternalError errorText="%s"`,
			ideaID, err)
		bark.RespondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	// create dog
	dog := &Dog{
		ID:           uuid.NewString(),
		CreationTime: time.Now(),
		IdeaID:       ideaID,
		ScheduleType: requestBody.ScheduleType,
		Schedule:     schedule,
	}

	dog, err = service.TasksClient.Register(r.Context(), dog)
	if err != nil {
		service.Logf(r, `action=RegisterDog dogID=%s ideaID=%s result=InternalError errorText="%s"`, err)
		bark.RespondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}
	service.Logf(r, `action=RegisterDog dogID=%s ideaID=%s result=OK`, dog.ID, dog.IdeaID)

	bark.RespondSuccess(w, http.StatusCreated, dog)
}

// GetDog is a handler for getting a dog.
func (service *Service) GetDog(w http.ResponseWriter, r *http.Request) {
	dogID := chi.URLParam(r, "dogID")

	// get dog from datastore
	dog, err := service.DogGetter.Get(r.Context(), dogID)
	switch status.Code(err) {
	case codes.OK:
		service.Logf(r, `action=GetDog dogID=%s result=OK`, dogID)
		bark.RespondSuccess(w, http.StatusOK, dog)
	case codes.NotFound:
		service.Logf(r, `action=GetDog dogID=%s result=NotFoundError`, dogID)
		bark.RespondError(w, http.StatusNotFound, fmt.Sprint("dog not found with ID: ", dogID))
	default:
		service.Logf(r, `action=GetDog dogID=%s result=InternalError errorText="%s"`,
			dogID, err)
		bark.RespondError(w, http.StatusInternalServerError, "internal error occurred")
	}
}

// DeleteDog is a handler for deleting a dog.
func (service *Service) DeleteDog(w http.ResponseWriter, r *http.Request) {
	dogID := chi.URLParam(r, "dogID")

	err := service.TasksClient.Unregister(r.Context(), dogID)
	if err != nil {
		service.Logf(r, `action=UnregisterDog dogID=%s result=InternalError errorText="%s"`,
			dogID, err)
		bark.RespondError(w, http.StatusInternalServerError, "internal error occurred")
		return
	}

	service.Logf(r, `action=UnregisterDog dogID=%s result=OK`, dogID)
	bark.RespondSuccess(w, http.StatusNoContent, nil)
}
