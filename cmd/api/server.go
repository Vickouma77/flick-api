package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (a *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.port),
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}

	a.logger.Info("starting server", "addr", srv.Addr, "env", a.config.env)

	return srv.ListenAndServe()
}
