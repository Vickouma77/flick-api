package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// TODO: Refactor later
const version = "1.0.0"

// Configuration settings
type config struct {
	port int
	env  string
}

// Dependencies for the HTTP handlers, helpers and middlewares.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	// Initialize flags
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	// Initialize a structure logger to write to standard out stream
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// An instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	// A custom HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start HTTP server
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
