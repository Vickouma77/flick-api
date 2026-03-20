package main

import (
	"fmt"
	"net/http"
)

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, a.config.env, version)

	w.Header().Set("Content-type", "application/json")

	w.Write([]byte(js))
}
