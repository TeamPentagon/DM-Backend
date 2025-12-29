// Package main is the entry point for the DM-Backend service.
// It sets up the HTTP server with Gin framework and configures API routes.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Config holds application configuration
type Config struct {
	Port           string
	AIEndpoint     string
	RequestTimeout time.Duration
}

// Request represents the AI generation request payload
type Request struct {
	MaxContextLength int     `json:"max_context_length"`
	MaxLength        int     `json:"max_length"`
	Prompt           string  `json:"prompt"`
	Quiet            bool    `json:"quiet"`
	RepPen           float64 `json:"rep_pen"`
	RepPenRange      int     `json:"rep_pen_range"`
	RepPenSlope      float64 `json:"rep_pen_slope"`
	Temperature      float64 `json:"temperature"`
	Tfs              int     `json:"tfs"`
	TopA             int     `json:"top_a"`
	TopK             int     `json:"top_k"`
	TopP             float64 `json:"top_p"`
	Typical          int     `json:"typical"`
}

// Response represents the API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DefaultConfig returns default configuration values
func DefaultConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	endpoint := os.Getenv("AI_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://herbal-pmc-allowance-cognitive.trycloudflare.com/api/v1/generate"
	}

	return Config{
		Port:           port,
		AIEndpoint:     endpoint,
		RequestTimeout: 30 * time.Second,
	}
}

// chatResponse handles chat requests and forwards them to the AI service
func chatResponse(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		prompt, exists := c.GetQuery("response")
		if !exists || prompt == "" {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Error:   "Missing 'response' query parameter",
			})
			return
		}

		log.Printf("Received chat request with prompt length: %d", len(prompt))

		data := Request{
			MaxContextLength: 2048,
			MaxLength:        100,
			Prompt:           prompt,
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

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling request: %v", err)
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to process request",
			})
			return
		}

		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: config.RequestTimeout,
		}

		resp, err := client.Post(config.AIEndpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error calling AI service: %v", err)
			c.JSON(http.StatusServiceUnavailable, Response{
				Success: false,
				Error:   "AI service unavailable",
			})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to read AI response",
			})
			return
		}

		log.Printf("AI response received, length: %d", len(body))

		c.JSON(http.StatusOK, Response{
			Success: true,
			Data:    json.RawMessage(body),
		})
	}
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"status": "healthy"},
	})
}

// setupRouter configures and returns the Gin router
func setupRouter(config Config) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/chat", chatResponse(config))
	}

	// Legacy route for backward compatibility
	router.POST("/chat", chatResponse(config))

	return router
}

func main() {
	config := DefaultConfig()

	log.Printf("Starting DM-Backend server on port %s", config.Port)
	log.Printf("AI Endpoint: %s", config.AIEndpoint)

	router := setupRouter(config)

	addr := fmt.Sprintf(":%s", config.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
