package main

import (
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/api"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/loadbalancer"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/scheduler"
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

	// Initialize the load balancer with the store.
	loadbalancer := loadbalancer.NewLoadBalancer(store)

	// Create a new API handler with the store and load balancer and set up the endpoints.
	handler := api.NewHandler(store, loadbalancer)

	// Start the scheduler in a separate goroutine to continuously schedule pending pods.
	scheduler := scheduler.NewScheduler(store)
	scheduler.Start()

	r.Get("/ping", handler.Ping)
	r.Post("/pods", handler.CreatePod)
	r.Get("/pods", handler.GetAllPods)
	r.Put("/pods/{id}/status", handler.UpdatePodStatus)
	r.Post("/services", handler.CreateService)
	r.Get("/services", handler.GetAllServices)
	r.Get("/services/{name}/next", handler.GetNextPod)
	r.Post("/nodes/register", handler.RegisterNode)
	r.Post("/nodes/{id}/heartbeat", handler.UpdateHeartbeat)
	r.Get("/nodes", handler.GetAllNodes)

	slog.Info("Server running on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}

}
