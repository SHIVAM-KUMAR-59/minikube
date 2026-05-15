package scheduler

import (
	"log/slog"
	"time"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

type Scheduler struct {
	store *store.Store
	nodeIDs []string
	counter   int
}

// NewScheduler creates a new Scheduler instance with the provided Store and hardcoded node IDs. The Scheduler is responsible for assigning Pods to nodes based on the available node IDs.
func NewScheduler(store *store.Store) *Scheduler {
	// Hardcoded node IDs for now
	return &Scheduler{store: store, nodeIDs: []string{"node1", "node2", "node3"}}
}

// launches a goroutine with a time.NewTicker that ticks every 5 seconds and calls a schedule() method each tick
func (s *Scheduler) Start() {
	slog.Info("Starting scheduler")
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			s.Schedule()
		}
	}()
}

// Schedule retrieves all Pods from the store and schedules any Pods that are in the Pending state. It assigns each pending Pod to a node in a round-robin fashion and updates the Pod's status to Scheduled. If there are any errors during the scheduling process, it logs the errors using slog.
func (s *Scheduler) Schedule() {
	pods, err := s.store.GetAllPods()
	if err != nil {
		slog.Error("Failed to get pods from store", "error", err)
		return
	}

	for _, pod := range pods {
		if pod.Status == store.StatusPending {
			// Pick a node in round robin fashion
			nodeID := s.nodeIDs[s.counter % len(s.nodeIDs)]
			s.counter++

			// Update the pod status to "Scheduled" and assign it to the selected node
			pod.Status = store.StatusScheduled
			pod.NodeID = nodeID

			err := s.store.UpdatePod(pod)
			if err != nil {
				slog.Error("Failed to update pod in store", "error", err)
				continue
			}

			slog.Info("Pod scheduled successfully", "pod_id", pod.ID, "node_id", nodeID)
		}
	}
}