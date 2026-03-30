package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				// If panic, Set "Connection: close" header
				w.Header().Set("Connection", "close")

				// Log custom error
				a.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	// client stores per-IP rate limit state.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Keep one limiter per client IP. A mutex protects the map because this
	// middleware is called concurrently by multiple goroutines.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Periodically remove stale clients so the map does not grow forever.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock during cleanup to synchronize with request-time reads/writes.
			mu.Lock()

			// Evict clients that have been inactive for more than three minutes.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse remote address and use only the IP as the per-client key.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		// Lock while accessing the shared clients map.
		mu.Lock()

		// First request from this IP: create a limiter with 2 requests/second and
		// a burst capacity of 4.
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// Update activity timestamp for stale-entry eviction.
		clients[ip].lastSeen = time.Now()

		// Allow consumes one token if available; otherwise reject with 429.
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			a.rateLimitExceededResponse(w, r)
			return
		}

		// Unlock before calling next to avoid serializing downstream handlers.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
