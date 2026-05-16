package store

import (
	"encoding/json"
	"log/slog"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Node represents a Kubernetes Node with relevant information for storage and retrieval.
type Node struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Status        string    `json:"status"`
}

// RegisterNode takes a Node struct, serializes the Node to JSON, and saves it in the "nodes" bucket of BoltDB using the node's ID as the key.
func (s *Store) RegisterNode(node Node) error {
	// Serialize the Node struct to JSON.
	nodeData, err := json.Marshal(node)
	if err != nil {
		slog.Error("Failed to serialize node", "error", err)
		return err
	}

	// Save the serialized node data to the "nodes" bucket with node.ID as the key.
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		return bucket.Put([]byte(node.ID), nodeData)
	})
	if err != nil {
		slog.Error("Failed to save node to BoltDB", "error", err)
		return err
	}

	return nil
}

// GetAllNodes retrieves all nodes from the "nodes" bucket in BoltDB, deserializes them from JSON, and returns a slice of Node structs.
func (s *Store) GetAllNodes() ([]Node, error) {
	var nodes []Node

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		return bucket.ForEach(func(k, v []byte) error {
			var node Node
			if err := json.Unmarshal(v, &node); err != nil {
				slog.Error("Failed to deserialize node", "error", err)
				return err
			}
			nodes = append(nodes, node)
			return nil
		})
	})

	return nodes, err
}

// UpdateNodeHeartbeat updates the LastHeartbeat field of a node with the current time. If the node is not found, it returns nil without an error.
func (s *Store) UpdateNodeHeartbeat(nodeID string) error {
	node, err := s.GetNodeByID(nodeID)
	if err != nil {
		slog.Error("Failed to retrieve node for heartbeat update", "error", err)
		return err
	}

	if node == nil {
		slog.Warn("Node not found for heartbeat update", "nodeID", nodeID)
		return nil // Node not found, return nil without error
	}

	node.LastHeartbeat = time.Now()

	return s.RegisterNode(*node)
}

// GetNodeByID retrieves a node from the "nodes" bucket in BoltDB by its ID, deserializes it from JSON, and returns a pointer to the Node struct. If the node is not found, it returns nil without an error.
func (s *Store) GetNodeByID(nodeID string) (*Node, error) {
	var node Node

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		nodeData := bucket.Get([]byte(nodeID))
		if nodeData == nil {
			return nil // Node not found, return nil without error
		}
		return json.Unmarshal(nodeData, &node)
	})

	if err != nil {
		slog.Error("Failed to retrieve node from BoltDB", "error", err)
		return nil, err
	}

	if node.ID == "" {
		return nil, nil // Node not found, return nil without error
	}

	return &node, nil
}

// DeleteNode deletes a node from the "nodes" bucket in BoltDB using the node ID as the key.
func (s *Store) DeleteNode(nodeID string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		return bucket.Delete([]byte(nodeID))
	})
	if err != nil {
		slog.Error("Failed to delete node from BoltDB", "error", err)
		return err
	}

	return nil
}
