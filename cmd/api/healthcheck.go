package main

import (
	"net/http"
)

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.env,
			"version":     version,
		},
	}

	err := a.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		a.serverError(w, r, err)
	}
}
