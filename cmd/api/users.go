package main

import (
	"errors"
	"net/http"

	"flick.io/internal/data"
	"flick.io/internal/validator"
)

func (a *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "user with this email already exists")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}
