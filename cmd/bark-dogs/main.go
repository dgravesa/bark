package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/firestore"
	"github.com/dgravesa/bark/pkg/bark"
	"github.com/dgravesa/bark/pkg/dog"
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
	r.Use(middleware.Timeout(15 * time.Second))

	// initialize dependencies
	firestoreClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		logger.Fatal(err)
	}
	tasksClient, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	doggoStore := &dog.DoggoFirestore{
		FirestoreClient: firestoreClient,
	}
	// initialize service
	service := dog.Service{
		Service: bark.Service{
			Name:   "bark-dogs",
			Logger: logger,
		},
		IdeaGetter: &bark.IdeaFirestore{
			FirestoreClient: firestoreClient,
		},
		DogGetter: doggoStore,
		TasksClient: &dog.Whisperer{
			QueueName:  os.Getenv("QUEUE_NAME"),
			TaskClient: tasksClient,
			DogStore:   doggoStore,
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
