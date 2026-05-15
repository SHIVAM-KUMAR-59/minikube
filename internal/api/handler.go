package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

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

// Ping handles the /ping endpoint, responding with a JSON message indicating that Minikube is running.
func (h *Handler) Ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	fmt.Fprintln(res, `{"status": "ok", "message": "Minikube is running"}`)
	slog.Info("Ping request received")
}

type CreatePodRequest struct {
	Name string `json:"name"`
	Image string `json:"image"`
}

func generatePodID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

// CreatePod handles the /pods endpoint for creating a new Pod. It decodes the request body to get the Pod name and image, validates the input, creates a new Pod struct, saves it to the store, and responds with a success message and the created Pod.
func (h *Handler) CreatePod(res http.ResponseWriter, req *http.Request) {
	// Take out name and image from the request body
	var createPodReq CreatePodRequest
	err := json.NewDecoder(req.Body).Decode(&createPodReq)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that name and image are provided
	if createPodReq.Name == "" || createPodReq.Image == "" {
		slog.Error("Name and image are required fields")
		http.Error(res, "Name and image are required fields", http.StatusBadRequest)
		return
	}

	// Create a new Pod struct with the provided name and image
	pod := store.Pod{
		ID: generatePodID(),
		Name: createPodReq.Name,
		Image: createPodReq.Image,
		Status: "Pending",
		NodeID: "",
	}

	// Save the pod to the store
	err = h.store.CreatePod(pod)
	if err != nil {
		slog.Error("Failed to create pod", "error", err)
		http.Error(res, "Failed to create pod", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(pod)
	slog.Info("Pod created successfully", "pod_id", pod.ID)
}

// GetAllPods handles the /pods endpoint for retrieving all Pods. It retrieves all pods from the store, encodes them as JSON, and responds with the list of pods.
func (h *Handler) GetAllPods(res http.ResponseWriter, req *http.Request) {
	pods, err := h.store.GetAllPods()
	if err != nil {
		slog.Error("Failed to retrieve pods", "error", err)
		http.Error(res, "Failed to retrieve pods", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	json.NewEncoder(res).Encode(pods)
	slog.Info("Retrieved all pods successfully", "pod_count", len(pods))
}