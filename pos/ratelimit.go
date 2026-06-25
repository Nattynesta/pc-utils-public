package main

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	rl := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
	go rl.cleanup(3 * time.Minute)
	return rl
}

func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
	}
	return limiter
}

func (rl *IPRateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, limiter := range rl.ips {
			if limiter.Tokens() == float64(rl.b) {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

var (
	globalLimiter = NewIPRateLimiter(200.0/60.0, 40)
	authLimiter   = NewIPRateLimiter(30.0/60.0, 10)
)

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := net.ParseIP(strings.TrimSpace(parts[0])); ip != nil {
			return ip.String()
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func withRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		path := r.URL.Path
		if strings.HasPrefix(path, "/static/") || strings.HasPrefix(path, "/audio/") {
			next.ServeHTTP(w, r)
			return
		}
		isAuth := path == "/login" || strings.HasPrefix(path, "/api/usuarios") || strings.HasPrefix(path, "/api/register") || strings.HasPrefix(path, "/register")

		if isAuth {
			if !authLimiter.GetLimiter(ip).Allow() {
				slog.Warn("rate limit exceeded (auth)", "ip", ip, "path", path)
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"rate_limit_exceeded","retry_after":60}`))
				return
			}
		} else {
			if !globalLimiter.GetLimiter(ip).Allow() {
				slog.Warn("rate limit exceeded (global)", "ip", ip, "path", path)
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"rate_limit_exceeded","retry_after":60}`))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
