package store

import (
	"encoding/json"
	"fmt"
	"log/slog"

	bolt "go.etcd.io/bbolt"
)

// Pod represents a Kubernetes Pod with relevant information for storage and retrieval.
type Pod struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Status string `json:"status"`
	NodeID string `json:"node_id"`
}

// CreatePod takes a Pod struct, serializes the Pod to JSON, and saves it in the "pods" bucket of BoltDB using the pod's ID as the key.
func (s *Store) CreatePod(pod Pod) error {
	// Serialize the Pod struct to JSON.
	podData, err := json.Marshal(pod)
	if err != nil {
		slog.Error("Failed to serialize pod", "error", err)
		return err
	}

	// Check if a pod with the same name already exists
	existingPod, err := s.GetPodByName(pod.Name)
	if err != nil {
		slog.Error("Failed to check for existing pod", "error", err)
		return err
	}
	if existingPod != nil {
		slog.Error("Pod with the same name already exists", "podName", pod.Name)
		return fmt.Errorf("pod with the same name already exists")
	}

	// Save the serialized pod data to the "pods" bucket with pod.ID as the key.
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		return bucket.Put([]byte(pod.ID), podData)
	})
	if err != nil {
		slog.Error("Failed to save pod to BoltDB", "error", err)
		return err
	}

	return nil
}

// GetAllPods retrieves all pods from the "pods" bucket in BoltDB, deserializes them from JSON, and returns a slice of Pod structs.
func (s *Store) GetAllPods() ([]Pod, error) {
	var pods []Pod
	
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		return bucket.ForEach(func(k, v []byte) error {
			var pod Pod
			if err := json.Unmarshal(v, &pod); err != nil {
				slog.Error("Failed to deserialize pod", "error", err)
				return err
			}
			pods = append(pods, pod)
			return nil
		})
	})
	if err != nil {
		slog.Error("Failed to retrieve pods from BoltDB", "error", err)
		return nil, err
	}

	return pods, nil
}

// GetPodByID retrieves a single Pod from the "pods" bucket in BoltDB using the provided podID. It deserializes the pod data from JSON and returns a pointer to the Pod struct. If the pod is not found, it returns nil without an error.
func (s *Store) GetPodByID(podID string) (*Pod, error) {
	var pod Pod

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		podData := bucket.Get([]byte(podID))
		if podData == nil {
			return nil // Pod not found, return nil without error
		}
		return json.Unmarshal(podData, &pod)
	})
	if err != nil {
		slog.Error("Failed to retrieve pod from BoltDB", "error", err)
		return nil, err
	}

	if pod.ID == "" {
		return nil, nil // Pod not found
	}

	return &pod, nil
}

// UpdatePod takes a Pod struct, serializes it to JSON, and updates the existing pod data in the "pods" bucket of BoltDB using the pod's ID as the key. If the pod does not exist, it will be created.
func (s *Store) UpdatePod(pod Pod) error {
	// Serialize the Pod struct to JSON.
	podData, err := json.Marshal(pod)
	if err != nil {
		slog.Error("Failed to serialize pod", "error", err)
		return err
	}

	// Update the existing pod data in the "pods" bucket with the new serialized data.
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		return bucket.Put([]byte(pod.ID), podData)
	})
	if err != nil {
		slog.Error("Failed to update pod in BoltDB", "error", err)
		return err
	}

	return nil
}

// DeletePod removes a pod from the "pods" bucket in BoltDB using the provided podID as the key. If the pod does not exist, it will simply return without an error.
func (s *Store) DeletePod(podID string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		return bucket.Delete([]byte(podID))
	})
	if err != nil {
		slog.Error("Failed to delete pod from BoltDB", "error", err)
		return err
	}
	
	return nil
}

// GetPodByName retrieves a single Pod from the "pods" bucket in BoltDB using the provided podName. It iterates through all pods, deserializes them from JSON, and returns a pointer to the Pod struct that matches the given name. If no pod with the specified name is found, it returns nil without an error.
func (s *Store) GetPodByName(podName string) (*Pod, error) {
	var pod *Pod
	
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("pods"))
		return bucket.ForEach(func(k, v []byte) error {
			var p Pod
			if err := json.Unmarshal(v, &p); err != nil {
				slog.Error("Failed to deserialize pod", "error", err)
				return err
			}
			if p.Name == podName {
				pod = &p
				return nil // Stop iterating once we find the pod
			}
			return nil
		})
	})
	if err != nil {
		slog.Error("Failed to retrieve pods from BoltDB", "error", err)
		return nil, err
	}

	return pod, nil
}