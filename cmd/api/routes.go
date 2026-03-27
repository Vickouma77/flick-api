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
	router.HandlerFunc(http.MethodGet, "/v1/movies", a.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", a.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", a.deleteMovieHandler)

	return a.recoverPanic(router)
}
