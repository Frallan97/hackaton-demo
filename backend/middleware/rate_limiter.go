package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiterEntry represents a single entry in the rate limiter
type RateLimiterEntry struct {
	Count     int
	FirstSeen time.Time
}

// RateLimiter implements a simple in-memory rate limiter using sliding window
type RateLimiter struct {
	entries map[string]*RateLimiterEntry
	mutex   sync.RWMutex
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a new rate limiter with specified limit and window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*RateLimiterEntry),
		limit:   limit,
		window:  window,
	}
	
	// Start cleanup goroutine to remove expired entries
	go rl.cleanup()
	
	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	
	entry, exists := rl.entries[ip]
	if !exists {
		// First request from this IP
		rl.entries[ip] = &RateLimiterEntry{
			Count:     1,
			FirstSeen: now,
		}
		return true
	}
	
	// Check if the window has expired
	if now.Sub(entry.FirstSeen) > rl.window {
		// Reset the counter for a new window
		entry.Count = 1
		entry.FirstSeen = now
		return true
	}
	
	// Check if limit is exceeded
	if entry.Count >= rl.limit {
		return false
	}
	
	// Increment counter
	entry.Count++
	return true
}

// cleanup removes expired entries from the map
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, entry := range rl.entries {
			if now.Sub(entry.FirstSeen) > rl.window {
				delete(rl.entries, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware creates middleware that applies rate limiting
func RateLimitMiddleware(rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP address from request
			ip := getClientIP(r)
			
			if !rateLimiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Rate limit exceeded. Please try again later."}`))
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (common in load balancers/proxies)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}