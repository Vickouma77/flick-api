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
	if err, ok := message.(error); ok {
		message = err.Error()
	}

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

// badRequestResponse returns a 400 Bad Request response with the provided error message.
func (a *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse returns a 422 Unprocessable Entity response containing validation errors.
func (a *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (a *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	a.errorResponse(w, r, http.StatusConflict, message)
}

func (a *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	a.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (a *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	a.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (a *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"

	a.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (a *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated access this resource"
	a.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (a *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	a.errorResponse(w, r, http.StatusForbidden, message)
}
