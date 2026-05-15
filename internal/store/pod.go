package store

// Pod represents a Kubernetes Pod with relevant information for storage and retrieval.
type Pod struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Status string `json:"status"`
	NodeID string `json:"node_id"`
}