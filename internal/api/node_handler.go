package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/go-chi/chi/v5"
)

// Handler provides methods to handle API requests related to nodes and other resources.
type RegisterNodeRequest struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

// RegisterNode handles the /register endpoint, allowing nodes to register themselves with the Minikube API. It decodes the incoming JSON request, creates a Node struct, and saves it to the BoltDB using the Store's RegisterNode method.
func (h *Handler) RegisterNode(res http.ResponseWriter, req *http.Request) {
	var registerNodeRequest RegisterNodeRequest

	// Decode the incoming JSON request body into the RegisterNodeRequest struct.
	if err := json.NewDecoder(req.Body).Decode(&registerNodeRequest); err != nil {
		slog.Error("Failed to decode register node request", "error", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a Node struct using the decoded request data and set the LastHeartbeat to the current time and Status to READY.
	node := store.Node{
		ID: registerNodeRequest.ID,
		Name: registerNodeRequest.Name,
		LastHeartbeat: time.Now(),
		Status: store.NodeStatusReady,
	}

	// Use the Store's RegisterNode method to save the node information to BoltDB. If there's an error, log it and return a 500 Internal Server Error response.
	if err := h.store.RegisterNode(node); err != nil {
		slog.Error("Failed to register node", "error", err)
		http.Error(res, "Failed to register node", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	fmt.Fprintln(res, `{"status": "ok", "message": "Node registered successfully"}`)
	slog.Info("Node registered successfully", "node_id", node.ID)
}

// UpdateHeartbeat handles the /{id}/heartbeat endpoint, allowing nodes to update their heartbeat information. It extracts the node ID from the URL, validates it, and uses the Store's UpdateNodeHeartbeat method to update the node's last heartbeat time in BoltDB. If successful, it responds with a success message in JSON format.
func (h *Handler) UpdateHeartbeat(res http.ResponseWriter, req *http.Request) {
	// Extract the ID from URL params
	nodeID := chi.URLParam(req, "id")

	if nodeID == "" {
		slog.Error("Node ID is required")
		http.Error(res, "Node ID is required", http.StatusBadRequest)
		return
	}

	// Update the node's heartbeat in the store
	if err := h.store.UpdateNodeHeartbeat(nodeID); err != nil {
		slog.Error("Failed to update node heartbeat", "error", err)
		http.Error(res, "Failed to update node heartbeat", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, `{"status": "ok", "message": "Node heartbeat updated successfully"}`)
	slog.Info("Node heartbeat updated successfully", "node_id", nodeID)
}
