package store

import (
	"log/slog"

	bolt "go.etcd.io/bbolt"
)

// Store provides methods to interact with BoltDB for storing and retrieving Pod information.
type Store struct {
	db *bolt.DB
}

// Opens a DB connection and creates a bucket called "pods", "services" and "nodes" if they don't exist, then returns a Store instance.
func NewStore(dbPath string) (*Store, error) {

	// Open the BoltDB database file. If it doesn't exist, it will be created.
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		slog.Error("Failed to open BoltDB", "error", err)
		return nil, err
	}

	// Create a bucket called "pods", "services" and "nodes" if they doesn't exist.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pods"))
		if err != nil {
			slog.Error("Failed to create bucket for pods", "error", err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("services"))
		if err != nil {
			slog.Error("Failed to create bucket for services", "error", err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			slog.Error("Failed to create bucket for nodes", "error", err)
			return err
		}
		
		return err
	})
	if err != nil {
		slog.Error("Failed to create bucket 'pods'", "error", err)
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close closes the BoltDB database connection when the Store is no longer needed.
func (s *Store) Close() error {
    return s.db.Close()
}