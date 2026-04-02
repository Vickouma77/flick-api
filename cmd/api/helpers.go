package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"flick.io/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

// readIDParam extracts the `id` route parameter and validates that it is a
// positive integer.
func (a *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// writeJSON marshals the provided value as JSON, applies any custom response
// headers, and writes the HTTP status code and body.
func (a *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a trailing newline to make JSON responses easier to read in logs and terminals.
	js = append(js, '\n')

	// Merge caller-supplied headers into the response header map.
	maps.Insert(w.Header(), maps.All(headers))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// readJSON decodes the request body into dest and converts common JSON decoding
// failures into client-friendly error messages.
func (a *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body into the target destination value.
	err := dec.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		// Map decoding failures to clearer validation-style messages.
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unkown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unkown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// readString returns the query-string value for key, or defaultValue if the
// parameter is missing or empty.
func (a *application) readString(qs url.Values, key, defaultValue string) string {
	// Extracting the value of a given key
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// readCSV returns a comma-separated query-string value as a slice. If the
// parameter is missing or empty, it returns defaultValue.
func (a *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// readInt parses an integer query-string value. If the parameter is missing,
// empty, or invalid, it records a validation error (for invalid values) and
// returns defaultValue.
func (a *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (a *application) background(fn func()) {
	// Incrementing waitGroup counter
	a.wg.Add(1)

	// Background goroutine
	go func() {
		// Use defer to decrement the WaitGroup counter before the goroutine returns.
		defer a.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				a.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
