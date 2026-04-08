package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	// Initialize a httprouter router instance
	router := httprouter.New()

	// Custom error handlers
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	// Register relevant methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", a.requirePermission("movies:read", a.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.requirePermission("movies:write", a.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.requirePermission("movies:read", a.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", a.requirePermission("movies:write", a.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", a.requirePermission("movies:write", a.deleteMovieHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return a.metrics(a.recoverPanic(a.enableCORS(a.rateLimit(a.authenticate(router)))))
}
