package worker

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Worker represents a worker node in the cluster, responsible for managing and executing tasks.
type Worker struct {
	dockerClient *client.Client
	store *store.Store
	nodeID string
}

// NewWorker creates a new Worker instance with the provided store and node ID.
func NewWorker(store *store.Store, nodeID string) (*Worker, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return nil, err
	}

	return &Worker{
		dockerClient: dockerClient,
		store: store,
		nodeID: nodeID,
	}, nil
}

// Start launches a goroutine that periodically calls the Reconcile method to check for scheduled pods and attempt to run them.
func (w *Worker) Start() {
	slog.Info("Worker started", "nodeID", w.nodeID)
	ticker := time.NewTicker(5 * time.Second)
	go func () {
		for range ticker.C {
			w.Reconcile()
		}
	}()
}

// Reconcile checks for any pods that are scheduled to run on this worker node and attempts to run them. It fetches all pods from the store, checks their status, and if a pod is in the Scheduled state and assigned to this worker's node ID, it calls the RunPod method to execute the pod. If there are any errors during this process, it logs the errors using slog.
func (w *Worker) Reconcile() {
	// Fetch all pods
	pods, err := w.store.GetAllPods()
	if err != nil {
		slog.Error("Failed to fetch pods from store", "error", err)
		return
	}

	// Iterate through the pods and check their status
	for _, pod := range pods {
		if pod.Status == store.StatusScheduled && pod.NodeID == w.nodeID {
			slog.Info("Pod is scheduled, attempting to run it", "podID", pod.ID)
			w.RunPod(pod)
		}
	}
}

// RunPod takes a pod as input and attempts to run it using the Docker client. It first pulls the required image, then creates a container based on that image, and finally starts the container. If any of these steps fail, it logs the error using slog. After successfully starting the container, it updates the pod's status to Running in the store.
func (w *Worker) RunPod (pod store.Pod) {
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
	}, nil, nil, nil, pod.Name)
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
	pod.Status = store.StatusRunning
	err = w.store.UpdatePod(pod)
	if err != nil {
		slog.Error("Failed to update pod status in store", "error", err)
		return
	}

	slog.Info("Pod is now running", "podID", pod.ID)
}