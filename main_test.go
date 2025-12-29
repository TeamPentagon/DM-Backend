// Package main provides unit tests for the main application entry point.
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestDefaultConfig tests the default configuration loading
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Port == "" {
		t.Error("Expected default port to be set")
	}

	if config.AIEndpoint == "" {
		t.Error("Expected default AI endpoint to be set")
	}

	if config.RequestTimeout <= 0 {
		t.Error("Expected request timeout to be positive")
	}
}

// TestSetupRouter tests router configuration
func TestSetupRouter(t *testing.T) {
	config := DefaultConfig()
	router := setupRouter(config)

	if router == nil {
		t.Fatal("Expected router to be created")
	}

	// Test that routes are registered
	routes := router.Routes()
	expectedRoutes := map[string]bool{
		"GET /health":        false,
		"POST /chat":         false,
		"POST /api/v1/chat":  false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		expectedRoutes[key] = true
	}

	for route, found := range expectedRoutes {
		if !found {
			t.Errorf("Expected route %s to be registered", route)
		}
	}
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	config := DefaultConfig()
	router := setupRouter(config)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if data["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", data["status"])
	}
}

// TestChatResponseMissingParameter tests chat endpoint with missing parameter
func TestChatResponseMissingParameter(t *testing.T) {
	config := DefaultConfig()
	router := setupRouter(config)

	req := httptest.NewRequest("POST", "/api/v1/chat", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error == "" {
		t.Error("Expected error message to be present")
	}
}

// TestChatResponseEmptyParameter tests chat endpoint with empty parameter
func TestChatResponseEmptyParameter(t *testing.T) {
	config := DefaultConfig()
	router := setupRouter(config)

	req := httptest.NewRequest("POST", "/api/v1/chat?response=", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestLegacyChatEndpoint tests the legacy /chat endpoint
func TestLegacyChatEndpoint(t *testing.T) {
	config := DefaultConfig()
	router := setupRouter(config)

	// Should work the same as /api/v1/chat
	req := httptest.NewRequest("POST", "/chat", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 400 for missing parameter (same behavior as new endpoint)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestResponseStructure tests the Response struct marshaling
func TestResponseStructure(t *testing.T) {
	tests := []struct {
		name     string
		response Response
	}{
		{
			name: "Success response",
			response: Response{
				Success: true,
				Data:    map[string]string{"key": "value"},
			},
		},
		{
			name: "Error response",
			response: Response{
				Success: false,
				Error:   "Something went wrong",
			},
		},
		{
			name: "Empty data response",
			response: Response{
				Success: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			var unmarshaled Response
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if unmarshaled.Success != tt.response.Success {
				t.Errorf("Success mismatch: got %v, want %v", unmarshaled.Success, tt.response.Success)
			}

			if unmarshaled.Error != tt.response.Error {
				t.Errorf("Error mismatch: got %v, want %v", unmarshaled.Error, tt.response.Error)
			}
		})
	}
}

// TestRequestStructure tests the Request struct
func TestRequestStructure(t *testing.T) {
	req := Request{
		MaxContextLength: 2048,
		MaxLength:        100,
		Prompt:           "Test prompt",
		Quiet:            false,
		RepPen:           1.1,
		RepPenRange:      256,
		RepPenSlope:      1,
		Temperature:      0.5,
		Tfs:              1,
		TopA:             0,
		TopK:             100,
		TopP:             0.9,
		Typical:          1,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var unmarshaled Request
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.Prompt != req.Prompt {
		t.Errorf("Prompt mismatch: got %v, want %v", unmarshaled.Prompt, req.Prompt)
	}

	if unmarshaled.Temperature != req.Temperature {
		t.Errorf("Temperature mismatch: got %v, want %v", unmarshaled.Temperature, req.Temperature)
	}

	if unmarshaled.MaxLength != req.MaxLength {
		t.Errorf("MaxLength mismatch: got %v, want %v", unmarshaled.MaxLength, req.MaxLength)
	}
}
