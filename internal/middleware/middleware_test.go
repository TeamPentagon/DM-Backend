// Package middleware provides middleware testing utilities.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestCORSMiddleware tests the CORS middleware
func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(CORS(DefaultCORSConfig()))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
	}{
		{
			name:           "Regular GET request",
			method:         "GET",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Preflight OPTIONS request",
			method:         "OPTIONS",
			origin:         "http://example.com",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check CORS headers
			if w.Header().Get("Access-Control-Allow-Origin") == "" && tt.method != "OPTIONS" {
				t.Error("Expected Access-Control-Allow-Origin header to be set")
			}
		})
	}
}

// TestLoggerMiddleware tests the logger middleware
func TestLoggerMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test?query=value", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestRateLimiter tests the rate limiter functionality
func TestRateLimiter(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 2,
		BurstSize:         3,
	}
	rl := NewRateLimiter(config)

	clientIP := "192.168.1.1"

	// First 3 requests should be allowed (burst)
	for i := 0; i < 3; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 4th request should be blocked
	if rl.Allow(clientIP) {
		t.Error("4th request should be blocked")
	}

	// Wait for token replenishment
	time.Sleep(1 * time.Second)

	// Should be allowed now
	if !rl.Allow(clientIP) {
		t.Error("Request after waiting should be allowed")
	}
}

// TestRateLimiterMiddleware tests the rate limiter middleware
func TestRateLimiterMiddleware(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 1,
		BurstSize:         2,
	}
	rl := NewRateLimiter(config)

	router := gin.New()
	router.Use(rl.RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// 3rd request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

// TestRecoveryMiddleware tests panic recovery
func TestRecoveryMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(Recovery())
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	// Should not panic
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthMiddleware tests authentication middleware
func TestAuthMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(Auth())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "With auth header",
			authHeader:     "Bearer token123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Without auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestRequestIDMiddleware tests request ID middleware
func TestRequestIDMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("RequestID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no request id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	tests := []struct {
		name            string
		providedID      string
		shouldUseProvided bool
	}{
		{
			name:            "With provided request ID",
			providedID:      "custom-request-id",
			shouldUseProvided: true,
		},
		{
			name:            "Without provided request ID",
			providedID:      "",
			shouldUseProvided: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.providedID != "" {
				req.Header.Set("X-Request-ID", tt.providedID)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			responseID := w.Header().Get("X-Request-ID")
			if responseID == "" {
				t.Error("Expected X-Request-ID header to be set")
			}

			if tt.shouldUseProvided && responseID != tt.providedID {
				t.Errorf("Expected request ID '%s', got '%s'", tt.providedID, responseID)
			}
		})
	}
}

// TestJoinStrings tests the string joining helper
func TestJoinStrings(t *testing.T) {
	tests := []struct {
		input    []string
		sep      string
		expected string
	}{
		{[]string{"a", "b", "c"}, ", ", "a, b, c"},
		{[]string{"a"}, ", ", "a"},
		{[]string{}, ", ", ""},
	}

	for _, tt := range tests {
		result := joinStrings(tt.input, tt.sep)
		if result != tt.expected {
			t.Errorf("joinStrings(%v, %s) = %s, want %s", tt.input, tt.sep, result, tt.expected)
		}
	}
}
