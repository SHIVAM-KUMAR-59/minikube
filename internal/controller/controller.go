package controller

import (
	"log/slog"
	"time"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

type Controller struct {
	db *store.Store
}

// NewController creates a new Controller instance with the provided Store. The Controller is responsible for resetting Pods whose assigned nodes are dead / unhealthy.
func NewController(db *store.Store) *Controller {
	return &Controller{
		db: db,
	}
}

// launches a goroutine with a time.NewTicker that ticks every 10 seconds and calls a reconcile() method each tick
func (c *Controller) Start() {
	slog.Info("Starting controller")
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			c.reconcile()
		}
	}()
}

// Reconcile checks for any nodes that are dead / unhealthy to run pods and attempts to re-assign the pods. It fetches all nodes from the store, checks their status, and if a node's last heartbeat is more than 15 seconds ago, it is declared dead. Any pods assigned to this node which were in RUNNING or SCHEDULED states are then reset back to PENDING for another node to pick them up If there are any errors during this process, it logs the errors using slog.
func (c *Controller) reconcile() {

	// Get all the nodes
	nodes, err := c.db.GetAllNodes()
	if err != nil {
		slog.Error("Failed to get all nodes", "error", err)
		return
	}

	currentTime := time.Now()

	for _, node := range nodes {
		if currentTime.Sub(node.LastHeartbeat) > 15*time.Second {
			// Node is unhealthy / heartbeat expired, mark it as NOT_READY
			err = c.db.UpdateNodeStatus(node.ID, store.NodeStatusNotReady)
			if err != nil {
				slog.Error("Failed to update node status to NOT_READY", "error", err)
			}

			// Find all pods assigned to that node with status RUNNING or SCHEDULED
			pods, err := c.db.GetPodsByNodeID(node.ID)
			if err != nil {
				slog.Error("Failed to fetch pods with node ID", "nodeID", node.ID, "error", err)
				return
			}

			// Reset them to PENDING so the scheduler picks them up again
			for _, pod := range pods {
				pod.Status = store.StatusPending
				c.db.UpdatePod(pod)
			}
		}
	}
}
