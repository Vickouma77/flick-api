package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"flick.io/internal/data"
	"flick.io/internal/validator"

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
		if a.config.limiter.enabled {
			// Parse remote address and use only the IP as the per-client key.
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				a.serverError(w, r, err)
				return
			}

			// Lock while accessing the shared clients map.
			mu.Lock()

			// First request from this IP: create a limiter from runtime config.
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst)}
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
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = a.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headersPart := strings.Split(authorizationHeader, " ")
		if len(headersPart) != 2 || headersPart[0] != "Bearer" {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headersPart[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := a.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				a.invalidAuthenticationTokenResponse(w, r)
			default:
				a.serverError(w, r, err)
			}
			return
		}

		r = a.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (a *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)

		if user.IsAnonymous() {
			a.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)

		if !user.Activated {
			a.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return a.requireAuthenticatedUser(fn)
}

func (a *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)

		permissions, err := a.models.Permissions.GetAllForUSer(int64(user.ID))
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		if !permissions.Include(code) {
			a.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

	return a.requireActivatedUser(fn)
}

func (a *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			if slices.Contains(a.config.cors.trustedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)

				if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
					w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

					w.WriteHeader(http.StatusOK)
					return 
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
