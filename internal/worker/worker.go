package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Worker represents a worker node in the cluster, responsible for managing and executing tasks.
type Worker struct {
	dockerClient *client.Client
	serverUrl    string
	nodeID       string
}

// NewWorker creates a new Worker instance with the provided store and node ID.
func NewWorker(serverUrl string, nodeID string) (*Worker, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return nil, err
	}

	return &Worker{
		dockerClient: dockerClient,
		serverUrl:    serverUrl,
		nodeID:       nodeID,
	}, nil
}

// Start launches a goroutine that periodically calls the Reconcile method to check for scheduled pods and attempt to run them.
func (w *Worker) Start() {
	slog.Info("Worker started", "nodeID", w.nodeID)

	reconcileTicker := time.NewTicker(5 * time.Second)
	heartbeatTicker := time.NewTicker(5 * time.Second)

	// Register the worker node
	time.Sleep(2 * time.Second) // Sleep for a short duration to ensure the API server is up and running before attempting to register the worker node.
	resp, err := http.Post(fmt.Sprintf("%s/nodes/register", w.serverUrl), "application/json", strings.NewReader(fmt.Sprintf(`{"id": "%s", "name": "%s"}`, w.nodeID, w.nodeID)))
	if err != nil {
		slog.Error("Failed to register worker node", "error", err)
		return
	}
	resp.Body.Close()

	go func() {
		for range reconcileTicker.C {
			w.Reconcile()
		}
	}()

	// Start a separate goroutine to send heartbeat signals to the API server at regular intervals.
	go func() {
		for range heartbeatTicker.C {
			// Send heartbeat to the API server
			http.Post(fmt.Sprintf("%s/nodes/%s/heartbeat", w.serverUrl, w.nodeID), "application/json", nil)
		}
	}()
}

// Reconcile checks for any pods that are scheduled to run on this worker node and attempts to run them. It fetches all pods from the store, checks their status, and if a pod is in the Scheduled state and assigned to this worker's node ID, it calls the RunPod method to execute the pod. If there are any errors during this process, it logs the errors using slog.
func (w *Worker) Reconcile() {
	// Fetch all pods
	response, err := http.Get(fmt.Sprintf("%s/pods", w.serverUrl))
	if err != nil {
		slog.Error("Failed to fetch pods from store", "error", err)
		return
	}
	defer response.Body.Close()

	var pods []store.Pod
	err = json.NewDecoder(response.Body).Decode(&pods)
	if err != nil {
		slog.Error("Failed to decode pods response", "error", err)
		return
	}

	// Log the number of pods fetched
	slog.Info("Fetched pods from store", "pod_count", len(pods))

	// Iterate through the pods and check their status
	for _, pod := range pods {
		if pod.Status == store.StatusScheduled && pod.NodeID == w.nodeID {
			slog.Info("Pod is scheduled, attempting to run it", "podID", pod.ID)
			w.RunPod(pod)
		}
	}
}

// RunPod takes a pod as input and attempts to run it using the Docker client. It first pulls the required image, then creates a container based on that image, and finally starts the container. If any of these steps fail, it logs the error using slog. After successfully starting the container, it updates the pod's status to Running in the store.
func (w *Worker) RunPod(pod store.Pod) {
	ctx := context.Background()

	// Pull the image
	reader, err := w.dockerClient.ImagePull(ctx, pod.Image, image.PullOptions{})
	if err != nil {
		slog.Error("Failed to pull image", "error", err)
		return
	}

	io.Copy(io.Discard, reader)
	reader.Close()

	// Create the container
	container, err := w.dockerClient.ContainerCreate(ctx, &dockerContainer.Config{
		Image: pod.Image,
	}, nil, nil, nil, pod.Name+"-"+pod.ID[:8])
	if err != nil {
		slog.Error("Failed to create container", "error", err)
		return
	}

	// Start the container
	err = w.dockerClient.ContainerStart(ctx, container.ID, dockerContainer.StartOptions{})
	if err != nil {
		slog.Error("Failed to start container", "error", err)
		return
	}

	// update pod status to StatusRunning in store
	reqBody := strings.NewReader(fmt.Sprintf(`{"status": "%s"}`, store.StatusRunning))
	httpReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/pods/%s/status", w.serverUrl, pod.ID), reqBody)
	if err != nil {
		slog.Error("Failed to create update request", "error", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		slog.Error("Failed to update pod status", "error", err)
		return
	}
	defer httpResp.Body.Close()

	slog.Info("Pod is now running", "podID", pod.ID)
}
