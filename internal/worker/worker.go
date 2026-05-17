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
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-chi/chi/v5"
)

// Worker represents a worker node in the cluster, responsible for managing and executing tasks.
type Worker struct {
	dockerClient *client.Client
	serverUrl    string
	nodeID       string
	port         string
}

// NewWorker creates a new Worker instance with the provided store and node ID.
func NewWorker(serverUrl string, nodeID string, port string) (*Worker, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return nil, err
	}

	return &Worker{
		dockerClient: dockerClient,
		serverUrl:    serverUrl,
		nodeID:       nodeID,
		port:         port,
	}, nil
}

func (w *Worker) StartHttpServer() {
	r := chi.NewRouter()

	// Endpoint to stream logs
	r.Get("/logs/{containerName}", func(res http.ResponseWriter, req *http.Request) {

		// Extract the container name
		containerName := chi.URLParam(req, "containerName")

		// Call docker to get the logs
		ctx := context.Background()
		logs, err := w.dockerClient.ContainerLogs(ctx, containerName, dockerContainer.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			slog.Error("Failed to get container logs", "container", containerName, "error", err)
			http.Error(res, "Failed to get logs", http.StatusInternalServerError)
			return
		}
		defer logs.Close()

		// Stream the logs
		res.Header().Set("Content-Type", "text/plain")
		stdcopy.StdCopy(res, res, logs)
	})

	// Goroutine to run the server as an unblocking operation
	go func() {
		if err := http.ListenAndServe(":"+w.port, r); err != nil {
			slog.Error("Failed to start worker HTTP server", "error", err)
		}
	}()
}

// Start launches a goroutine that periodically calls the Reconcile method to check for scheduled pods and attempt to run them.
func (w *Worker) Start() {
	slog.Info("Worker started", "nodeID", w.nodeID)

	reconcileTicker := time.NewTicker(5 * time.Second)
	heartbeatTicker := time.NewTicker(5 * time.Second)

	// Register the worker node
	time.Sleep(2 * time.Second) // Sleep for a short duration to ensure the API server is up and running before attempting to register the worker node.
	resp, err := http.Post(fmt.Sprintf("%s/nodes/register", w.serverUrl), "application/json", strings.NewReader(fmt.Sprintf(`{"id": "%s", "name": "%s", "address": "http://localhost:%s"}`, w.nodeID, w.nodeID, w.port)))
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

// Reconcile checks all pods assigned to this worker node and ensures they are in the correct state. For pods in SCHEDULED state, it calls runPod to start the container. For pods in RUNNING state, it inspects the container via Docker — if the container is missing or stopped, it removes it and restarts it automatically. Any errors during this process are logged using slog.
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
			// Run the pod
			slog.Info(
				"Pod is scheduled, attempting to run it",
				"podID", pod.ID,
			)
			w.runPod(pod)
		}
	}

	// Finds all pods with status RUNNING and NodeID == w.nodeID
	for _, pod := range pods {
		if pod.Status == store.StatusRunning && pod.NodeID == w.nodeID {
			containerName := pod.Name + "-" + pod.ID[:8]
			ctx := context.Background()

			// For each one, calls ContainerInspect using the container name
			containerStatus, err := w.dockerClient.ContainerInspect(
				ctx,
				containerName,
			)

			// If the container state is exited or the container doesn't exist — call runPod(pod) again to restart it
			if err != nil {
				slog.Warn("Container missing, restarting pod", "podID", pod.ID, "container", containerName, "error", err)
				w.removeContainer(containerName)
				w.runPod(pod)
				continue
			}

			// Container exists but is not running
			if !containerStatus.State.Running {
				slog.Warn("Container stopped, restarting pod", "podID", pod.ID, "container", containerName, "status", containerStatus.State.Status)
				w.removeContainer(containerName)
				w.runPod(pod)
			}
		}
	}
}

// runPod takes a pod as input and attempts to run it using the Docker client. It first pulls the required image, then creates a container based on that image, and finally starts the container. If any of these steps fail, it logs the error using slog. After successfully starting the container, it updates the pod's status to Running in the store.
func (w *Worker) runPod(pod store.Pod) {
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

// removeContainer takes a container name and removes that container from docker to avoid conflicts
func (w *Worker) removeContainer(containerName string) {
	ctx := context.Background()
	err := w.dockerClient.ContainerRemove(ctx, containerName, dockerContainer.RemoveOptions{
		Force: true,
	})
	if err != nil {
		slog.Warn("Failed to remove container", "container", containerName, "error", err)
	}
}
