// Package middleware provides HTTP middleware components for the DM-Backend service.
// It includes common middleware such as authentication, logging, rate limiting, and CORS handling.
package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration options
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns a permissive CORS configuration for development
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// CORS returns a CORS middleware handler
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range config.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders, ", "))
		
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Max-Age", formatDuration(config.MaxAge))

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logger returns a request logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		log.Printf("[%d] %s %s %s %v %s",
			status,
			method,
			path,
			query,
			latency,
			clientIP,
		)
	}
}

// RateLimiter configuration
type RateLimiterConfig struct {
	RequestsPerSecond int
	BurstSize         int
}

// RateLimiter implements a simple rate limiter using token bucket algorithm
type RateLimiter struct {
	mu       sync.Mutex
	tokens   map[string]int
	lastTime map[string]time.Time
	config   RateLimiterConfig
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		tokens:   make(map[string]int),
		lastTime: make(map[string]time.Time),
		config:   config,
	}
}

// RateLimitMiddleware returns a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rl.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if a request from the given client IP is allowed
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lastTime, exists := rl.lastTime[clientIP]

	if !exists {
		// New client, initialize with full tokens
		rl.tokens[clientIP] = rl.config.BurstSize - 1
		rl.lastTime[clientIP] = now
		return true
	}

	// Calculate tokens to add based on elapsed time
	elapsed := now.Sub(lastTime)
	tokensToAdd := int(elapsed.Seconds()) * rl.config.RequestsPerSecond
	currentTokens := rl.tokens[clientIP] + tokensToAdd

	// Cap at burst size
	if currentTokens > rl.config.BurstSize {
		currentTokens = rl.config.BurstSize
	}

	if currentTokens <= 0 {
		return false
	}

	rl.tokens[clientIP] = currentTokens - 1
	rl.lastTime[clientIP] = now
	return true
}

// Recovery returns a recovery middleware that handles panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Auth returns an authentication middleware
// Currently checks for the presence of an Authorization header
// TODO: Implement proper JWT validation
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// TODO: Validate JWT token here
		// For now, just check if header exists

		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Helper functions

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func formatDuration(d time.Duration) string {
	return string(rune(int(d.Seconds())))
}

func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
