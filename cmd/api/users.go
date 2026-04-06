package main

import (
	"errors"
	"net/http"
	"time"

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

	err = a.models.Permissions.AddForUSer(int64(user.ID), "movies:read")
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	token, err := a.models.Token.New(int64(user.ID), 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	// Send the welcome email in the background so SMTP issues don't block signup.
	// userCopy := *user
	a.background(func() {
		data := map[string]any{
			"activationToken": token.PlainText,
			"userID":          user.ID,
		}

		err = a.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			a.logger.Error(err.Error())
		}
	})

	err = a.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}

func (a *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := a.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	user.Activated = true

	err = a.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.editConflictResponse(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	err = a.models.Token.DeleteAllForUser(data.ScopeActivation, int64(user.ID))
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}
