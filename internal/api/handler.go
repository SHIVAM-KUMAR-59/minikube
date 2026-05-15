package api

import (
	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/google/uuid"
)

// Handler provides methods to handle HTTP requests related to Pod operations.
type Handler struct {
	store *store.Store
}

// NewHandler creates a new Handler instance with the provided Store.
func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func generateRandomID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

