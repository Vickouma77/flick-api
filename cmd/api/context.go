package main

import (
	"context"
	"net/http"

	"flick.io/internal/data"
)

// Custom contextKey type
type contextKey string

const userContextKey = contextKey("user")

func (a *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (a *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in the request context")
	}

	return user
}
