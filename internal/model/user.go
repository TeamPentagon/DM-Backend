// Package model provides data models and database operations for the DM-Backend application.
// It includes user management functionality with CRUD operations using LevelDB and Protocol Buffers.
package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/TeamPentagon/DM-Backend/internal/database"
	"google.golang.org/protobuf/proto"
)

// Common errors for user operations
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDatabaseOpen    = errors.New("failed to open database")
	ErrSerialize       = errors.New("failed to serialize user data")
	ErrDeserialize     = errors.New("failed to deserialize user data")
	ErrDatabaseWrite   = errors.New("failed to write to database")
	ErrDatabaseRead    = errors.New("failed to read from database")
	ErrDatabaseDelete  = errors.New("failed to delete from database")
	ErrInvalidUsername = errors.New("username cannot be empty")
)

// SaveUserData persists the user data to the LevelDB database.
// It uses the PersonId as the key for storage.
// Returns an error if the database operation fails.
func (s *User) SaveUserData() error {
	if s.PersonId == "" {
		return ErrInvalidUsername
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	bytes, err := proto.Marshal(s)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		return fmt.Errorf("%w: %v", ErrSerialize, err)
	}

	err = db.Put([]byte(s.PersonId), bytes, nil)
	if err != nil {
		log.Printf("Error writing to database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseWrite, err)
	}

	return nil
}

// GetUserData retrieves user data from the database using the provided username.
// The retrieved data is unmarshaled into the User struct receiver.
// Returns an error if the user is not found or if database operations fail.
func (g *User) GetUserData(userName string) error {
	if userName == "" {
		return ErrInvalidUsername
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	data, err := db.Get([]byte(userName), nil)
	if err != nil {
		log.Printf("Error reading from database for user %s: %v", userName, err)
		return fmt.Errorf("%w: %v", ErrDatabaseRead, err)
	}

	err = proto.Unmarshal(data, g)
	if err != nil {
		log.Printf("Error unmarshaling user data: %v", err)
		return fmt.Errorf("%w: %v", ErrDeserialize, err)
	}

	return nil
}

// UpdateUserData updates the user data in the database.
// It retrieves the existing user by username, updates with new data from the receiver,
// and writes the updated data back to the database.
// Returns an error if the user is not found or if database operations fail.
func (u *User) UpdateUserData(userName string) error {
	if userName == "" {
		return ErrInvalidUsername
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	// Check if user exists
	_, err = db.Get([]byte(userName), nil)
	if err != nil {
		log.Printf("Error reading from database for user %s: %v", userName, err)
		return fmt.Errorf("%w: %v", ErrDatabaseRead, err)
	}

	// Marshal updated user data
	bytes, err := proto.Marshal(u)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		return fmt.Errorf("%w: %v", ErrSerialize, err)
	}

	// Write updated data back to database
	err = db.Put([]byte(userName), bytes, nil)
	if err != nil {
		log.Printf("Error writing to database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseWrite, err)
	}

	return nil
}

// DeleteUserData removes the user data from the database.
// Returns an error if the user is not found or if the delete operation fails.
func (d *User) DeleteUserData(userName string) error {
	if userName == "" {
		return ErrInvalidUsername
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	// Check if user exists before attempting delete
	_, err = db.Get([]byte(userName), nil)
	if err != nil {
		log.Printf("Error reading from database for user %s: %v", userName, err)
		return fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	// Delete the user data
	err = db.Delete([]byte(userName), nil)
	if err != nil {
		log.Printf("Error deleting from database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseDelete, err)
	}

	return nil
}

// UserExists checks if a user exists in the database.
// Returns true if the user exists, false otherwise.
func UserExists(userName string) (bool, error) {
	if userName == "" {
		return false, ErrInvalidUsername
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	_, err = db.Get([]byte(userName), nil)
	if err != nil {
		return false, nil // User doesn't exist
	}

	return true, nil
}
