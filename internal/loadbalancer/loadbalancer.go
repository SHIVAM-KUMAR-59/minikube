package loadbalancer

import (
	"log/slog"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
)

// LoadBalancer is responsible for distributing incoming traffic to the appropriate pods based on the service's load balancing strategy (e.g., round-robin).
type LoadBalancer struct {
	db *store.Store
	currentCounter map[string]int // Key is service ID, value is the current counter for round-robin
}

// NewLoadBalancer initializes a new LoadBalancer instance with the provided store. It also initializes the currentCounter map to keep track of the round-robin counters for each service.
func NewLoadBalancer(db *store.Store) *LoadBalancer {
	return &LoadBalancer{
		db: db,
		currentCounter: make(map[string]int),
	}
}

// GetNextPodForService retrieves the next pod for the given service name using a round-robin strategy. It first retrieves the service from the store, gets the list of associated pod IDs, and then selects the next pod ID based on the current counter for that service. Finally, it retrieves and returns the corresponding Pod struct from the store.
func (lb *LoadBalancer) GetNextPodForService(serviceName string) (*store.Pod, error) {

	// Retrieve the service with the given name
	service, err := lb.db.GetServiceByName(serviceName)
	if err != nil {
		slog.Error("Failed to retrieve service from store", "error", err)
		return nil, err
	}

	// Get the list of pod IDs associated with the service
	podIDs := service.Pods
	if len(podIDs) == 0 {
		return nil, nil // No pods available for this service
	}

	// Get the current counter for this service, defaulting to 0 if not set
	counter := lb.currentCounter[service.ID]

	// Select the next pod ID using round-robin strategy
	selectedPodID := podIDs[counter%len(podIDs)]

	// Increment the counter for the next request
	lb.currentCounter[service.ID] = (counter + 1) % len(podIDs)

	// Retrieve the pod with the selected pod ID
	pod, err := lb.db.GetPodByID(selectedPodID)
	if err != nil {
		slog.Error("Failed to retrieve pod from store", "error", err)
		return nil, err
	}

	return pod, nil
}