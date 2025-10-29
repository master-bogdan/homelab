package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// internal token bucket
type bucket struct {
	mu       sync.Mutex
	tokens   int
	max      int
	interval time.Duration
	last     time.Time
}

func newBucket(max int, refillInterval time.Duration) *bucket {
	return &bucket{
		tokens:   max,
		max:      max,
		interval: refillInterval,
		last:     time.Now(),
	}
}

func (b *bucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.last)

	if elapsed >= b.interval {
		b.tokens = b.max
		b.last = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// exported RateLimiter wrapper with per-IP buckets
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	max      int
	interval time.Duration
}

func NewRateLimiter(max int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucket),
		max:      max,
		interval: interval,
	}
}

func clientIP(r *http.Request) string {
	// X-Forwarded-For can contain multiple IPs: client, proxies...
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) getBucket(ip string) *bucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[ip]
	if !ok {
		b = newBucket(rl.max, rl.interval)
		rl.buckets[ip] = b
	}
	return b
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}

		ip := clientIP(r)
		if !rl.getBucket(ip).allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
