package api

import (
	"github.com/SHIVAM-KUMAR-59/minikube/internal/loadbalancer"
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/google/uuid"
)

// Handler provides methods to handle HTTP requests related to Pod operations.
type Handler struct {
	store *store.Store
	loadBalancer *loadbalancer.LoadBalancer
}

// NewHandler creates a new Handler instance with the provided Store.
func NewHandler(store *store.Store, loadBalancer *loadbalancer.LoadBalancer) *Handler {
	return &Handler{store: store, loadBalancer: loadBalancer}
}

// Utility function to generate a random ID for services and pods.
func generateRandomID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

