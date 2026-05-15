package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

// Handler provides methods to handle HTTP requests related to Pod operations.
type Handler struct {
	store *store.Store
}

// NewHandler creates a new Handler instance with the provided Store.
func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

// Ping handles the /ping endpoint, responding with a JSON message indicating that Minikube is running.
func (h *Handler) Ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	fmt.Fprintln(res, `{"status": "ok", "message": "Minikube is running"}`)
	slog.Info("Ping request received")
}