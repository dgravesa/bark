package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/dgravesa/bark/pkg/bark"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var r *chi.Mux

// HealthResponse is a basic health response type.
type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

func init() {
	r = chi.NewRouter()

	// initialize middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := &HealthResponse{
			Healthy: true,
		}
		bark.RespondSuccess(w, http.StatusOK, response)
	})
}

func main() {
	portNum := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *portNum)
	http.ListenAndServe(addr, r)
}
