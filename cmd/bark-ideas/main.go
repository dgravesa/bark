package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dgravesa/bark/pkg/bark"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const projectID = "itsthecloudforyourcloud"

var (
	r      *chi.Mux
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	r = chi.NewRouter()

	// initialize middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// initialize service
	firestoreClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		logger.Fatal(err)
	}
	service := bark.IdeaService{
		Service: bark.Service{
			Name:   "bark-ideas",
			Logger: logger,
		},
		IdeaStore: &bark.IdeaFirestore{
			FirestoreClient: firestoreClient,
		},
	}

	service.RegisterRoutes(r)
}

func main() {
	portNum := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *portNum)
	logger.Printf("listenPort=%d", *portNum)
	http.ListenAndServe(addr, r)
}
