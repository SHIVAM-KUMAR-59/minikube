package store

import (
	"encoding/json"
	"log/slog"

	bolt "go.etcd.io/bbolt"
)

// Store provides methods to interact with BoltDB for storing and retrieving Pod information.
type Store struct {
	db *bolt.DB
}

// Opens a DB connection and creates a bucket called "pods" if it doesn't exist, then returns a Store instance.
func NewStore(dbPath string) (*Store, error) {

	// Open the BoltDB database file. If it doesn't exist, it will be created.
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		slog.Error("Failed to open BoltDB", "error", err)
		return nil, err
	}

	// Create a bucket called "pods" if it doesn't exist.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pods"))
		return err
	})
	if err != nil {
		slog.Error("Failed to create bucket 'pods'", "error", err)
		return nil, err
	}

	return &Store{db: db}, nil
}

// CreatePod takes a Pod struct, serializes the Pod to JSON, and saves it in the "pods" bucket of BoltDB using the pod's ID as the key.
func (s *Store) CreatePod(pod Pod) error {
	// Serialize the Pod struct to JSON.
	podData, err := json.Marshal(pod)
	if err != nil {
		slog.Error("Failed to serialize pod", "error", err)
		return err
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

// Close closes the BoltDB database connection when the Store is no longer needed.
func (s *Store) Close() error {
    return s.db.Close()
}