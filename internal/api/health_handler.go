package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

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
		slog.Error("Failed to fetch pods", "error", err)

		http.Error(res, "Failed to fetch pods", http.StatusInternalServerError)

		return
	}

	// Fetch all nodes
	nodes, err := h.store.GetAllNodes()
	if err != nil {
		slog.Error("Failed to fetch nodes", "error", err)

		http.Error(res, "Failed to fetch nodes", http.StatusInternalServerError)

		return
	}

	// Fetch all services
	services, err := h.store.GetAllServices()
	if err != nil {
		slog.Error("Failed to fetch services", "error", err)

		http.Error(res, "Failed to fetch services", http.StatusInternalServerError)

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
		slog.Error("Failed to encode cluster health response", "error", err)
	}
}
