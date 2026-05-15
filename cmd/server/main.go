package main

import (
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/api"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

func main() {
	mux := http.NewServeMux()

	// Initialize the BoltDB store and create a new Handler for API endpoints.
	store, err := store.NewStore("minikube.db")
	if err != nil {
		slog.Error("Error creating store", "error", err)
		return
	}
	defer store.Close()

	// Create a new API handler with the store and set up the endpoints.
	handler := api.NewHandler(store)

	mux.HandleFunc("/ping", handler.Ping)

	slog.Info("Server running on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}