package main

import (
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

	router.HandlerFunc(http.MethodGet, "/v1/movies", a.requireActivatedUser(a.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.requireActivatedUser(a.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.requireActivatedUser(a.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", a.requireActivatedUser(a.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", a.requireActivatedUser(a.deleteMovieHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))
}
