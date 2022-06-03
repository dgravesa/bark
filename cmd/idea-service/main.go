package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/bark/pkg/bark"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const projectName = "ItsTheCloudForYourCloud"

var r *chi.Mux

func init() {
	r = chi.NewRouter()

	// initialize middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// initialize service
	datastoreClient, err := datastore.NewClient(context.Background(), projectName)
	if err != nil {
		log.Fatal(err)
	}
	service := IdeaService{
		IdeaStore: &bark.IdeaDatastore{
			DatastoreClient: datastoreClient,
		},
	}

	service.RegisterRoutes(r)
}

func main() {
	portNum := flag.Int("port", 8345, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *portNum)
	log.Printf("listening on port %d...", *portNum)
	http.ListenAndServe(addr, r)
}
