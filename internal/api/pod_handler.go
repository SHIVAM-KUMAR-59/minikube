package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

type CreatePodRequest struct {
	Name string `json:"name"`
	Image string `json:"image"`
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
		ID: generateRandomID(),
		Name: createPodReq.Name,
		Image: createPodReq.Image,
		Status: store.StatusPending,
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