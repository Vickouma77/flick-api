package main

import (
	"fmt"
	"net/http"
)

// logError records request context alongside an error for troubleshooting.
func (a *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

// errorResponse writes a standardized JSON error payload with the given status code.
func (a *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := a.writeJSON(w, status, env, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverError returns a generic 500 response and logs the underlying error.
func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"

	a.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse returns a 404 when no route/resource matches the request.
func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"

	a.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowed returns a 405 when the route exists but does not support the request method.
func (a *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	a.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
