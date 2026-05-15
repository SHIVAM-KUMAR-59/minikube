package api

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Ping handles the /ping endpoint, responding with a JSON message indicating that Minikube is running.
func (h *Handler) Ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	fmt.Fprintln(res, `{"status": "ok", "message": "Minikube is running"}`)
	slog.Info("Ping request received")
}