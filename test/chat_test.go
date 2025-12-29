// Package test provides integration tests for the DM-Backend chat functionality.
package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TeamPentagon/DM-Backend/internal/model"
)

// setupChatTestEnvironment creates a clean test environment for chat tests
func setupChatTestEnvironment(t *testing.T) func() {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	testDir := filepath.Join(os.TempDir(), "dm-backend-chat-test")
	os.MkdirAll(testDir, 0755)
	os.Chdir(testDir)

	return func() {
		os.Chdir(originalDir)
		os.RemoveAll(testDir)
	}
}

// createTestMessage creates a message with test data
func createTestMessage() *model.Message {
	return &model.Message{
		AiId:           "ai_001",
		UserId:         "user_001",
		Content:        "Hello, this is a test message",
		Timestamp:      time.Now().Unix(),
		ConversationId: "conv_001",
		MsgId:          "msg_001",
	}
}

// TestSaveMessage tests saving a message
func TestSaveMessage(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	msg := createTestMessage()

	err := msg.SaveMessage()
	if err != nil {
		t.Errorf("SaveMessage() returned an error: %v", err)
	}
}

// TestSaveMessageEmptyMsgId tests validation
func TestSaveMessageEmptyMsgId(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	msg := &model.Message{
		AiId:    "ai_001",
		UserId:  "user_001",
		Content: "Test",
		MsgId:   "", // Empty - should fail
	}

	err := msg.SaveMessage()
	if err == nil {
		t.Error("SaveMessage() should return error for empty MsgId")
	}
}

// TestSaveChatHistory tests saving chat history
func TestSaveChatHistory(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	history := &model.ChatHistory{
		ConversationId: "conv_001",
		Messages: []*model.Message{
			createTestMessage(),
		},
	}

	err := history.SaveChatHistory()
	if err != nil {
		t.Errorf("SaveChatHistory() returned an error: %v", err)
	}
}

// TestSaveChatHistoryEmptyConvId tests validation
func TestSaveChatHistoryEmptyConvId(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	history := &model.ChatHistory{
		ConversationId: "", // Empty - should fail
		Messages:       []*model.Message{},
	}

	err := history.SaveChatHistory()
	if err == nil {
		t.Error("SaveChatHistory() should return error for empty ConversationId")
	}
}

// TestGetChatHistory tests retrieving chat history
func TestGetChatHistory(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	// First save a chat history
	originalHistory := &model.ChatHistory{
		ConversationId: "conv_get_test",
		Messages: []*model.Message{
			{
				MsgId:          "msg_001",
				Content:        "First message",
				ConversationId: "conv_get_test",
			},
			{
				MsgId:          "msg_002",
				Content:        "Second message",
				ConversationId: "conv_get_test",
			},
		},
	}

	err := originalHistory.SaveChatHistory()
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Now retrieve it
	retrievedHistory := &model.ChatHistory{
		ConversationId: "conv_get_test",
	}

	err = retrievedHistory.GetChatHistory()
	if err != nil {
		t.Errorf("GetChatHistory() returned an error: %v", err)
	}

	if len(retrievedHistory.Messages) != len(originalHistory.Messages) {
		t.Errorf("Expected %d messages, got %d", len(originalHistory.Messages), len(retrievedHistory.Messages))
	}
}

// TestGetChatHistoryNotFound tests retrieving non-existent history
func TestGetChatHistoryNotFound(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	history := &model.ChatHistory{
		ConversationId: "nonexistent_conv",
	}

	err := history.GetChatHistory()
	if err == nil {
		t.Error("GetChatHistory() should return error for non-existent conversation")
	}
}

// TestAddMessageToHistory tests adding messages to history
func TestAddMessageToHistory(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	// Create a new conversation
	history := &model.ChatHistory{
		ConversationId: "conv_add_test",
	}

	// Add first message (should create new history)
	msg1 := &model.Message{
		MsgId:          "msg_add_001",
		Content:        "First message",
		ConversationId: "conv_add_test",
	}

	err := history.AddMessageToHistory(msg1)
	if err != nil {
		t.Errorf("AddMessageToHistory() returned an error: %v", err)
	}

	// Add second message
	msg2 := &model.Message{
		MsgId:          "msg_add_002",
		Content:        "Second message",
		ConversationId: "conv_add_test",
	}

	err = history.AddMessageToHistory(msg2)
	if err != nil {
		t.Errorf("AddMessageToHistory() second call returned an error: %v", err)
	}

	// Verify messages were added
	verifyHistory := &model.ChatHistory{
		ConversationId: "conv_add_test",
	}
	err = verifyHistory.GetChatHistory()
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}

	if len(verifyHistory.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(verifyHistory.Messages))
	}
}

// TestChatLifecycle tests the complete chat lifecycle
func TestChatLifecycle(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	conversationId := "conv_lifecycle"

	// Create conversation with initial message
	history := &model.ChatHistory{
		ConversationId: conversationId,
	}

	// Add messages
	messages := []struct {
		msgId   string
		content string
		fromAI  bool
	}{
		{"msg_lc_001", "Hello!", false},
		{"msg_lc_002", "Hi! How can I help you?", true},
		{"msg_lc_003", "What's the weather like?", false},
		{"msg_lc_004", "I cannot access real-time weather data.", true},
	}

	for _, m := range messages {
		msg := &model.Message{
			MsgId:          m.msgId,
			Content:        m.content,
			ConversationId: conversationId,
			Timestamp:      time.Now().Unix(),
		}
		if m.fromAI {
			msg.AiId = "assistant"
		} else {
			msg.UserId = "user_001"
		}

		err := history.AddMessageToHistory(msg)
		if err != nil {
			t.Fatalf("Failed to add message %s: %v", m.msgId, err)
		}
	}

	// Retrieve and verify
	finalHistory := &model.ChatHistory{
		ConversationId: conversationId,
	}
	err := finalHistory.GetChatHistory()
	if err != nil {
		t.Fatalf("Failed to get final history: %v", err)
	}

	if len(finalHistory.Messages) != len(messages) {
		t.Errorf("Expected %d messages, got %d", len(messages), len(finalHistory.Messages))
	}

	// Verify message order and content
	for i, msg := range finalHistory.Messages {
		if msg.GetContent() != messages[i].content {
			t.Errorf("Message %d: expected content '%s', got '%s'", i, messages[i].content, msg.GetContent())
		}
	}
}

// TestMultipleConversations tests handling multiple conversations
func TestMultipleConversations(t *testing.T) {
	cleanup := setupChatTestEnvironment(t)
	defer cleanup()

	conversations := []string{"conv_multi_001", "conv_multi_002", "conv_multi_003"}

	// Create multiple conversations
	for i, convId := range conversations {
		history := &model.ChatHistory{
			ConversationId: convId,
		}

		msg := &model.Message{
			MsgId:          convId + "_msg_001",
			Content:        convId + " message",
			ConversationId: convId,
			Timestamp:      time.Now().Unix() + int64(i),
		}

		err := history.AddMessageToHistory(msg)
		if err != nil {
			t.Fatalf("Failed to create conversation %s: %v", convId, err)
		}
	}

	// Verify each conversation is independent
	for _, convId := range conversations {
		history := &model.ChatHistory{
			ConversationId: convId,
		}

		err := history.GetChatHistory()
		if err != nil {
			t.Errorf("Failed to get conversation %s: %v", convId, err)
			continue
		}

		if len(history.Messages) != 1 {
			t.Errorf("Conversation %s: expected 1 message, got %d", convId, len(history.Messages))
		}

		expectedContent := convId + " message"
		if history.Messages[0].GetContent() != expectedContent {
			t.Errorf("Conversation %s: expected content '%s', got '%s'", convId, expectedContent, history.Messages[0].GetContent())
		}
	}
}
