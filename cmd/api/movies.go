package main

import (
	"errors"
	"fmt"
	"net/http"

	"flick.io/internal/data"
	"flick.io/internal/validator"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Movies.Insert(movie)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = a.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		a.serverError(w, r, err)
	}
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	movie, err := a.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}

func (a *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	movie, err := a.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	// Holds expected data from client
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genre   []string      `json:"genres"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genre != nil {
		movie.Genres = input.Genre
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Movies.Update(movie)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}

func (a *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}
