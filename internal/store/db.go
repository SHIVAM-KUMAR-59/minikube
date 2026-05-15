package store

import (
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

	// Create a bucket called "pods" and "services" if they doesn't exist.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pods"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("services"))
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