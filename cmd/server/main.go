package main

import (
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/api"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Initialize the BoltDB store and create a new Handler for API endpoints.
	store, err := store.NewStore("minikube.db")
	if err != nil {
		slog.Error("Error creating store", "error", err)
		return
	}
	defer store.Close()

	// Create a new API handler with the store and set up the endpoints.
	handler := api.NewHandler(store)

	r.Get("/ping", handler.Ping)
	r.Post("/pods", handler.CreatePod)
	r.Get("/pods", handler.GetAllPods)

	slog.Info("Server running on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}