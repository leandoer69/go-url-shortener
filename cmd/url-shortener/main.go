package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang-url-shortener/internal/config"
	"golang-url-shortener/internal/constants"
	"golang-url-shortener/internal/http-server/handlers/redirect"
	"golang-url-shortener/internal/http-server/handlers/url/delete"
	"golang-url-shortener/internal/http-server/handlers/url/save"
	"golang-url-shortener/internal/http-server/handlers/url/update"
	"golang-url-shortener/internal/http-server/middleware/logger"
	"golang-url-shortener/internal/lib/logger/sl"
	"golang-url-shortener/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.Login: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
		r.Put("/", update.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start a server")
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case constants.EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case constants.EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case constants.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
