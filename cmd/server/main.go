package main

import (
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/api"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/controller"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/loadbalancer"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/scheduler"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	// Start the Controller in a separate goroutine to continuously reassign dead node's pods.
	controller := controller.NewController(store)
	controller.Start()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Routes for API endpoints
	r.Get("/ping", handler.Ping)

	// Pod endpoints
	r.Post("/pods", handler.CreatePod)
	r.Get("/pods", handler.GetAllPods)
	r.Get("/pods/{podName}/logs", handler.GetPodLogs)
	r.Put("/pods/{id}/status", handler.UpdatePodStatus)
	r.Get("/pods/{podName}", handler.GetPodByName)
	r.Delete("/pods/{id}", handler.DeletePod)

	// Service endpoints
	r.Post("/services", handler.CreateService)
	r.Get("/services", handler.GetAllServices)
	r.Get("/services/{name}/next", handler.GetNextPod)
	r.Delete("/services/{id}", handler.DeleteService)

	// Node endpoints
	r.Post("/nodes/register", handler.RegisterNode)
	r.Post("/nodes/{id}/heartbeat", handler.UpdateHeartbeat)
	r.Get("/nodes", handler.GetAllNodes)
	r.Delete("/nodes/{id}", handler.DeleteNode)

	slog.Info("Server running on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}

}
