package store

import (
	"encoding/json"
	"log/slog"

	bolt "go.etcd.io/bbolt"
)

// Service represents a Kubernetes Service with relevant information for storage and retrieval.
type Service struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Pods []string `json:"pods"`
	Port string `json:"port"`
}

// CreateService takes a Service struct, serializes the Service to JSON, and saves it in the "services" bucket of BoltDB using the service's ID as the key.
func (s *Store) CreateService(service Service) error {
	// Serialize the Service struct to JSON.
	serviceData, err := json.Marshal(service)
	if err != nil {
		slog.Error("Failed to serialize service", "error", err)
		return err
	}

	// Save the serialized service data to the "services" bucket with service.ID as the key.
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("services"))
		return bucket.Put([]byte(service.ID), serviceData)
	})
	if err != nil {
		slog.Error("Failed to save service to BoltDB", "error", err)
		return err
	}

	return nil
}

// GetAllServices retrieves all services from the "services" bucket in BoltDB, deserializes them from JSON, and returns a slice of Service structs.
func (s *Store) GetAllServices() ([]Service, error) {
	var services []Service
	
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("services"))
		return bucket.ForEach(func(k, v []byte) error {
			var service Service
			if err := json.Unmarshal(v, &service); err != nil {
				slog.Error("Failed to deserialize service", "error", err)
				return err
			}
			services = append(services, service)
			return nil
		})
	})
	if err != nil {
		slog.Error("Failed to retrieve services from BoltDB", "error", err)
		return nil, err
	}

	return services, nil
}

// GetServiceByName retrieves a service by its name from the "services" bucket in BoltDB, deserializes it from JSON, and returns a pointer to the Service struct.
func (s *Store) GetServiceByName(name string) (*Service, error) {
	var service *Service
	
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("services"))
		return bucket.ForEach(func(k, v []byte) error {
			var s Service
			if err := json.Unmarshal(v, &s); err != nil {
				slog.Error("Failed to deserialize service", "error", err)
				return err
			}
			if s.Name == name {
				service = &s
				return nil // Stop iterating once we find the service
			}
			return nil
		})
	})
	if err != nil {
		slog.Error("Failed to retrieve services from BoltDB", "error", err)
		return nil, err
	}

	if service == nil {
		slog.Info("Service not found", "name", name)
		return nil, nil // Return nil if service is not found
	}

	return service, nil
}