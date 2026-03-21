package main

import (
	"encoding/json"
	"net/http"
)

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": a.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}
