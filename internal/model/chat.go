// Package model provides chat-related data models and database operations.
// It includes message handling functionality with CRUD operations using LevelDB and Protocol Buffers.
package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/TeamPentagon/DM-Backend/internal/database"
	"google.golang.org/protobuf/proto"
)

// Common errors for chat operations
var (
	ErrMessageNotFound     = errors.New("message not found")
	ErrChatHistoryNotFound = errors.New("chat history not found")
	ErrInvalidMessageID    = errors.New("message ID cannot be empty")
	ErrInvalidConvID       = errors.New("conversation ID cannot be empty")
)

// SaveMessage persists a message to the LevelDB database.
// Uses the message ID as the key with a "chat_" prefix.
// Returns an error if the database operation fails.
func (msg *Message) SaveMessage() error {
	if msg.GetMsgId() == "" {
		return ErrInvalidMessageID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	byteData, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("%w: %v", ErrSerialize, err)
	}

	err = db.Put([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), byteData, nil)
	if err != nil {
		log.Printf("Error writing message to database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseWrite, err)
	}

	return nil
}

// SaveChatHistory persists a chat history to the LevelDB database.
// Uses the conversation ID as the key with a "chat_" prefix.
// Returns an error if the database operation fails.
func (msg *ChatHistory) SaveChatHistory() error {
	if msg.GetConversationId() == "" {
		return ErrInvalidConvID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	byteData, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling chat history: %v", err)
		return fmt.Errorf("%w: %v", ErrSerialize, err)
	}

	err = db.Put([]byte(fmt.Sprintf("history_%s", msg.GetConversationId())), byteData, nil)
	if err != nil {
		log.Printf("Error writing chat history to database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseWrite, err)
	}

	return nil
}

// Delete removes a message from the database.
// Returns an error if the message is not found or if the delete operation fails.
func (msg *Message) Delete() error {
	if msg.GetMsgId() == "" {
		return ErrInvalidMessageID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	key := []byte(fmt.Sprintf("chat_%s", msg.GetMsgId()))

	// Check if message exists
	_, err = db.Get(key, nil)
	if err != nil {
		log.Printf("Error: message %s not found: %v", msg.GetMsgId(), err)
		return fmt.Errorf("%w: %v", ErrMessageNotFound, err)
	}

	err = db.Delete(key, nil)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseDelete, err)
	}

	return nil
}

// Update modifies an existing message in the database.
// Returns an error if the message is not found or if the update operation fails.
func (msg *Message) Update() error {
	if msg.GetMsgId() == "" {
		return ErrInvalidMessageID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	key := []byte(fmt.Sprintf("chat_%s", msg.GetMsgId()))

	// Check if message exists
	_, err = db.Get(key, nil)
	if err != nil {
		log.Printf("Error: message %s not found: %v", msg.GetMsgId(), err)
		return fmt.Errorf("%w: %v", ErrMessageNotFound, err)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("%w: %v", ErrSerialize, err)
	}

	// Update by putting new data (no need to delete first in LevelDB)
	err = db.Put(key, data, nil)
	if err != nil {
		log.Printf("Error updating message: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseWrite, err)
	}

	return nil
}

// Get retrieves a message from the database using its message ID.
// The retrieved data is unmarshaled into the Message struct receiver.
// Returns an error if the message is not found or if database operations fail.
func (msg *Message) Get() error {
	if msg.GetMsgId() == "" {
		return ErrInvalidMessageID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	data, err := db.Get([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), nil)
	if err != nil {
		log.Printf("Error reading message %s: %v", msg.GetMsgId(), err)
		return fmt.Errorf("%w: %v", ErrMessageNotFound, err)
	}

	err = proto.Unmarshal(data, msg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return fmt.Errorf("%w: %v", ErrDeserialize, err)
	}

	return nil
}

// GetChatHistory retrieves a chat history from the database using its conversation ID.
// Returns an error if the history is not found or if database operations fail.
func (ch *ChatHistory) GetChatHistory() error {
	if ch.GetConversationId() == "" {
		return ErrInvalidConvID
	}

	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}
	defer db.Close()

	data, err := db.Get([]byte(fmt.Sprintf("history_%s", ch.GetConversationId())), nil)
	if err != nil {
		log.Printf("Error reading chat history %s: %v", ch.GetConversationId(), err)
		return fmt.Errorf("%w: %v", ErrChatHistoryNotFound, err)
	}

	err = proto.Unmarshal(data, ch)
	if err != nil {
		log.Printf("Error unmarshaling chat history: %v", err)
		return fmt.Errorf("%w: %v", ErrDeserialize, err)
	}

	return nil
}

// AddMessageToHistory adds a message to an existing chat history.
// If the history doesn't exist, it creates a new one.
func (ch *ChatHistory) AddMessageToHistory(msg *Message) error {
	if ch.GetConversationId() == "" {
		return ErrInvalidConvID
	}
	if msg == nil {
		return ErrInvalidMessageID
	}

	// Try to get existing history
	err := ch.GetChatHistory()
	if err != nil {
		// Create new history if not found
		ch.Messages = []*Message{}
	}

	// Append the new message
	ch.Messages = append(ch.Messages, msg)

	// Save the updated history
	return ch.SaveChatHistory()
}
