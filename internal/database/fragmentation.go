// Package database provides fragmentation management for distributed data across shards.
// It uses LevelDB to maintain a global schema that maps row IDs to their respective shards.
package database

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

const (
	// GlobalSchemaPath is the default path for the fragmentation schema database
	GlobalSchemaPath = "GlobalSchema"
	// DefaultBaseDir is the default base directory for database files
	DefaultBaseDir = "Database/"
)

// Fragmentation errors
var (
	ErrShardNotFound      = errors.New("shard not found for the given key")
	ErrFragmentationWrite = errors.New("failed to write fragmentation entry")
	ErrFragmentationRead  = errors.New("failed to read fragmentation entry")
	ErrShardConversion    = errors.New("failed to convert shard to int64")
	ErrEmptyKey           = errors.New("key cannot be empty")
)

// FragmentationAdd adds one or more fragment entries to the global schema.
// Each key is mapped to the specified shard number.
//
// Parameters:
//   - shard: The shard number to associate with the keys
//   - keys: One or more keys to add to the fragmentation schema
//
// Returns an error if any database operation fails.
func FragmentationAdd(shard int64, keys ...string) error {
	if len(keys) == 0 {
		return ErrEmptyKey
	}

	ldb, err := CreateLevelDBDatabase(DefaultBaseDir, GlobalSchemaPath, "/")
	if err != nil {
		log.Printf("Error creating LevelDB database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer ldb.Close()

	shardBytes := []byte(fmt.Sprintf("%d", shard))
	for _, key := range keys {
		if key == "" {
			log.Printf("Warning: skipping empty key")
			continue
		}

		err = ldb.Put([]byte(key), shardBytes, nil)
		if err != nil {
			log.Printf("Error adding fragmentation entry for key %s: %v", key, err)
			return fmt.Errorf("%w for key %s: %v", ErrFragmentationWrite, key, err)
		}
		log.Printf("Added fragmentation: key=%s -> shard=%d", key, shard)
	}

	return nil
}

// FragmentationRemove removes a fragment entry from the global schema.
//
// Parameters:
//   - key: The key to remove from the fragmentation schema
//
// Returns an error if the database operation fails.
func FragmentationRemove(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	ldb, err := CreateLevelDBDatabase(DefaultBaseDir, GlobalSchemaPath, "/")
	if err != nil {
		log.Printf("Error creating LevelDB database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer ldb.Close()

	// Check if key exists before attempting to delete
	_, err = ldb.Get([]byte(key), nil)
	if err != nil {
		log.Printf("Key %s not found in fragmentation schema: %v", key, err)
		return fmt.Errorf("%w: %v", ErrShardNotFound, err)
	}

	err = ldb.Delete([]byte(key), nil)
	if err != nil {
		log.Printf("Error removing fragmentation entry for key %s: %v", key, err)
		return fmt.Errorf("%w: %v", ErrFragmentationWrite, err)
	}

	log.Printf("Removed fragmentation entry for key: %s", key)
	return nil
}

// FragmentationGet retrieves the shard number for a given row ID.
//
// Parameters:
//   - rowID: The row ID to look up in the fragmentation schema
//
// Returns the shard number and any error encountered.
func FragmentationGet(rowID string) (int64, error) {
	if rowID == "" {
		return -1, ErrEmptyKey
	}

	ldb, err := CreateLevelDBDatabase(DefaultBaseDir, GlobalSchemaPath, "/")
	if err != nil {
		log.Printf("Error creating LevelDB database: %v", err)
		return -1, fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer ldb.Close()

	shard, err := ldb.Get([]byte(rowID), nil)
	if err != nil {
		log.Printf("Error getting shard for rowID %s: %v", rowID, err)
		return -1, fmt.Errorf("%w: %v", ErrShardNotFound, err)
	}

	n, err := strconv.ParseInt(string(shard), 10, 64)
	if err != nil {
		log.Printf("Error converting shard to int64: %v", err)
		return -1, fmt.Errorf("%w: %v", ErrShardConversion, err)
	}

	return n, nil
}

// FragmentationExists checks if a key exists in the fragmentation schema.
//
// Parameters:
//   - key: The key to check
//
// Returns true if the key exists, false otherwise.
func FragmentationExists(key string) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}

	ldb, err := CreateLevelDBDatabase(DefaultBaseDir, GlobalSchemaPath, "/")
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer ldb.Close()

	_, err = ldb.Get([]byte(key), nil)
	if err != nil {
		return false, nil // Key doesn't exist
	}

	return true, nil
}

// FragmentationUpdate updates the shard assignment for an existing key.
//
// Parameters:
//   - key: The key to update
//   - newShard: The new shard number to assign
//
// Returns an error if the key doesn't exist or if the update fails.
func FragmentationUpdate(key string, newShard int64) error {
	if key == "" {
		return ErrEmptyKey
	}

	exists, err := FragmentationExists(key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: key %s", ErrShardNotFound, key)
	}

	ldb, err := CreateLevelDBDatabase(DefaultBaseDir, GlobalSchemaPath, "/")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer ldb.Close()

	err = ldb.Put([]byte(key), []byte(fmt.Sprintf("%d", newShard)), nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFragmentationWrite, err)
	}

	log.Printf("Updated fragmentation: key=%s -> shard=%d", key, newShard)
	return nil
}
