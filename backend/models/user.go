package models

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// User model holds user data.
type User struct {
	Username        string
	PasswordHash    string
	Role            string
	Allowed         bool
	AccessExpiresAt time.Time
}

// Users is an in-memory store for registered users.
var Users = map[string]User{}

// OpenAIRateLimit is the global rate limit (requests per minute).
var OpenAIRateLimit = 1

// RateLimiter tracks last request times per user.
type RateLimiter struct {
	mu           sync.Mutex
	lastRequests map[string]time.Time
}

// RateLimiterInstance is the singleton instance.
var RateLimiterInstance = NewRateLimiter()

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{lastRequests: make(map[string]time.Time)}
}

// Allow returns true if the user is allowed a new request.
func (r *RateLimiter) Allow(username string, rateLimit int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	interval := time.Minute / time.Duration(rateLimit)
	if last, ok := r.lastRequests[username]; ok {
		if now.Sub(last) < interval {
			return false
		}
	}
	r.lastRequests[username] = now
	return true
}

// OpenAIRequest represents a request to the OpenAI API.
type OpenAIRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// OpenAIResponse represents a response from OpenAI.
type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// RequestLog holds a record of a generate-post request.
type RequestLog struct {
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	requestLogsMu sync.Mutex
	requestLogs   []RequestLog
)

// LogRequest adds a new log entry.
func LogRequest(username string) {
	requestLogsMu.Lock()
	defer requestLogsMu.Unlock()
	requestLogs = append(requestLogs, RequestLog{
		Username:  username,
		Timestamp: time.Now(),
	})
}

// RequestLogs returns a copy of the logged requests.
func RequestLogs() []RequestLog {
	requestLogsMu.Lock()
	defer requestLogsMu.Unlock()
	cpy := make([]RequestLog, len(requestLogs))
	copy(cpy, requestLogs)
	return cpy
}

// PreCreateAdmin creates a default admin user.
func PreCreateAdmin() {
	if _, exists := Users["admin"]; !exists {
		// Use a pre-generated bcrypt hash for "adminpass" (for example purposes).
		Users["admin"] = User{
			Username:        "admin",
			PasswordHash:    "$2a$10$EixZaYVK1fsbw1ZfbX3OXe.PYp7OWmv09N5N6SOK6cc9J6cV67o7e",
			Role:            "admin",
			Allowed:         true,
			AccessExpiresAt: time.Now().Add(10080 * time.Minute),
		}
	}
}

// GetenvOrFail returns an environment variable or panics.
func GetenvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return val
}
