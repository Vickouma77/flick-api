package main

import (
	"errors"
	"net/http"
	"time"

	"flick.io/internal/data"
	"flick.io/internal/validator"
)

func (a *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := a.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentialsResponse(w, r)
			return
		default:
			a.serverError(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if !match {
		a.invalidCredentialsResponse(w, r)
		return
	}

	token, err := a.models.Token.New(int64(user.ID), 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}
