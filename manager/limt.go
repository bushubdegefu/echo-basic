package manager

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr

		if count, ok := rl.requests[ip]; ok {
			rl.requests[ip] = count + 1
			if count >= rl.limit {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		} else {
			rl.requests[ip] = 1
		}

		go func() {
			time.Sleep(rl.window)
			rl.mu.Lock()
			defer rl.mu.Unlock()
			delete(rl.requests, ip)
		}()

		next.ServeHTTP(w, r)
	})
}

func RatelimmterMain() {
	limit := 100
	window := time.Minute
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	rl := NewRateLimiter(limit, window)
	http.Handle("/", rl.Limit(handler))

	http.ListenAndServe(":8080", nil)
}
