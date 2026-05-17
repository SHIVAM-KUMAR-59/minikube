package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePodRequest struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
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

	replicas := createPodReq.Replicas
	if replicas == 0 {
		replicas = 1
	}

	var createdPods []store.Pod
	for i := 1; i <= replicas; i++ {
		podName := fmt.Sprintf("%s-%d", createPodReq.Name, i)

		// Create a new Pod struct with the provided name and image
		pod := store.Pod{
			ID:       generateRandomID(),
			Name:     podName,
			Image:    createPodReq.Image,
			Status:   store.StatusPending,
			NodeID:   "",
			Replicas: replicas,
		}

		// Save the pod to the store
		err = h.store.CreatePod(pod)
		if err != nil {
			slog.Error("Failed to create pod", "error", err)
			http.Error(res, "Failed to create pod", http.StatusInternalServerError)
			return
		}

		createdPods = append(createdPods, pod)
		slog.Info("Pod created successfully", "pod_id", pod.ID, "pod_name", pod.Name)
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(createdPods)
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

// GetPodByName handles the /pods/{podName} endpoint for retrieving one pod by its name. It retrieves the corresponding pod from the store, encodes it as JSON, and responds with the pod.
func (h *Handler) GetPodByName(res http.ResponseWriter, req *http.Request) {
	// Extract the pod name
	podName := chi.URLParam(req, "podName")

	if podName == "" {
		slog.Error("Pod name is required for deletion")
		http.Error(res, "Pod name is required", http.StatusBadRequest)
		return
	}

	// Fetch the pod by Name
	pod, err := h.store.GetPodByName(podName)
	if err != nil {
		slog.Error("Failed to fetch pod", "pod_name", podName, "error", err)
		http.Error(res, "Pod with this ID was not found", http.StatusNotFound)
		return
	}

	// Check if pod exists
	if pod == nil {
		slog.Warn("Pod not found", "pod_name", podName)
		http.Error(res, "Pod not found", http.StatusNotFound)
		return
	}

	slog.Info("Pod fetched successfully", "pod_id", pod.ID)
	res.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(res).Encode(pod); err != nil {
		slog.Error("Failed to encode pod response", "error", err)
	}
}

// UpdatePodStatus is a helper function that updates the status of a pod with the given pod ID. It retrieves the pod from the store, updates its status, and saves the updated pod back to the store. It also logs the status update operation.
func (h *Handler) UpdatePodStatus(res http.ResponseWriter, req *http.Request) {
	// Extract pod ID from the URL path
	podID := chi.URLParam(req, "id")

	if podID == "" {
		slog.Error("Pod ID is required")
		http.Error(res, "Pod ID is required", http.StatusBadRequest)
		return
	}

	// Extract the new status from the request body
	var updateReq struct {
		Status string `json:"status"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	pod, err := h.store.GetPodByID(podID)
	if err != nil {
		slog.Error("Failed to get pod for status update", "pod_id", podID, "error", err)
		return
	}

	pod.Status = updateReq.Status

	err = h.store.UpdatePod(*pod)
	if err != nil {
		slog.Error("Failed to update pod status", "pod_id", podID, "error", err)
		return
	}

	slog.Info("Pod status updated successfully", "pod_id", podID, "new_status", updateReq.Status)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, `{"status": "ok", "message": "Pod status updated successfully"}`)
}

// DeletePod handles the /pods/{id} endpoint for deleting a Pod. It extracts the pod ID from the URL, validates it, deletes the pod from the store, and responds with a success message if the deletion is successful.
func (h *Handler) DeletePod(res http.ResponseWriter, req *http.Request) {
	// Extract pod ID from the URL path
	podID := chi.URLParam(req, "id")

	if podID == "" {
		slog.Error("Pod ID is required for deletion")
		http.Error(res, "Pod ID is required", http.StatusBadRequest)
		return
	}

	err := h.store.DeletePod(podID)
	if err != nil {
		slog.Error("Failed to delete pod", "pod_id", podID, "error", err)
		http.Error(res, "Failed to delete pod", http.StatusInternalServerError)
		return
	}

	slog.Info("Pod deleted successfully", "pod_id", podID)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, `{"status": "ok", "message": "Pod deleted successfully"}`)
}

// GetPodLogs handles the /pods/{name}/logs endpoint for fetching the logs for a pod. It extracts the pod name from the URL, then fetches the corresponding pod from the store, get the node handling that pod and makes a GET request to node.Address/logs/{containerName}, then streams the response back to the client.
func (h *Handler) GetPodLogs(res http.ResponseWriter, req *http.Request) {
	// Extract the pod name
	podName := chi.URLParam(req, "podName")

	if podName == "" {
		slog.Error("Pod name is required for deletion")
		http.Error(res, "Pod name is required", http.StatusBadRequest)
		return
	}

	// Fetch the pod by Name
	pod, err := h.store.GetPodByName(podName)
	if err != nil {
		slog.Error("Failed to fetch pod", "pod_name", podName, "error", err)
		http.Error(res, "Pod with this ID was not found", http.StatusNotFound)
		return
	}

	// Fetch the node handling the pod
	node, err := h.store.GetNodeByID(pod.NodeID)
	if err != nil {
		slog.Error("Failed to fetch node", "node_id", pod.NodeID, "error", err)
		http.Error(res, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	// Build container name
	containerName := pod.Name + "-" + pod.ID[:8]

	// Make GET request to worker's log endpoint
	logsResp, err := http.Get(fmt.Sprintf("%s/logs/%s", node.Address, containerName))
	if err != nil {
		slog.Error("Failed to fetch logs from worker", "node_id", node.ID, "error", err)
		http.Error(res, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}
	defer logsResp.Body.Close()

	// Stream the response back to the client
	res.Header().Set("Content-Type", "text/plain")
	io.Copy(res, logsResp.Body)
	slog.Info("Logs streamed successfully", "pod_name", podName)
}

type ClusterHealthResponse struct {
	TotalPods     int    `json:"total_pods"`
	RunningPods   int    `json:"running_pods"`
	PendingPods   int    `json:"pending_pods"`
	TotalNodes    int    `json:"total_nodes"`
	ReadyNodes    int    `json:"ready_nodes"`
	TotalServices int    `json:"total_services"`
	ClusterHealth string `json:"cluster_health"`
}

func (h *Handler) GetClusterHealth(
	res http.ResponseWriter,
	req *http.Request,
) {

	// Fetch all pods
	pods, err := h.store.GetAllPods()
	if err != nil {
		slog.Error(
			"Failed to fetch pods",
			"error", err,
		)

		http.Error(
			res,
			"Failed to fetch pods",
			http.StatusInternalServerError,
		)

		return
	}

	// Fetch all nodes
	nodes, err := h.store.GetAllNodes()
	if err != nil {
		slog.Error(
			"Failed to fetch nodes",
			"error", err,
		)

		http.Error(
			res,
			"Failed to fetch nodes",
			http.StatusInternalServerError,
		)

		return
	}

	// Fetch all services
	services, err := h.store.GetAllServices()
	if err != nil {
		slog.Error(
			"Failed to fetch services",
			"error", err,
		)

		http.Error(
			res,
			"Failed to fetch services",
			http.StatusInternalServerError,
		)

		return
	}

	runningPods := 0
	pendingPods := 0

	for _, pod := range pods {

		switch pod.Status {
		case store.StatusRunning:
			runningPods++

		case store.StatusPending, store.StatusScheduled:
			pendingPods++
		}
	}

	readyNodes := 0

	for _, node := range nodes {
		if node.Status == store.NodeStatusReady {
			readyNodes++
		}
	}

	clusterHealth := "HEALTHY"

	if readyNodes == 0 {
		clusterHealth = "CRITICAL"
	} else if pendingPods > 0 {
		clusterHealth = "DEGRADED"
	}

	response := ClusterHealthResponse{
		TotalPods:     len(pods),
		RunningPods:   runningPods,
		PendingPods:   pendingPods,
		TotalNodes:    len(nodes),
		ReadyNodes:    readyNodes,
		TotalServices: len(services),
		ClusterHealth: clusterHealth,
	}

	res.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(res).Encode(response); err != nil {
		slog.Error(
			"Failed to encode cluster health response",
			"error", err,
		)
	}
}
