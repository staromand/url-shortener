package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	verbose := false

	for _, arg := range os.Args {
		if arg == "-v" {
			verbose = true
		}
	}

	log := logger.SetupLogger(cfg.Env, verbose)
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	dbStorage, err := sqlite.MigrateNew(cfg.StoragePath)
	if err != nil {
		log.Error("An error occurred while starting the application", sl.Err(err))
		os.Exit(1)
	}

	router.Post("/url", save.New(log, dbStorage))
	log.Info("starting http-server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
