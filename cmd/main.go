package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/save"
	mwLogger "github.com/leyl1ne/rest-api-parser/pkg/http-server/middleware/logger"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogpretty"
	"github.com/leyl1ne/rest-api-parser/pkg/storage/psql"
)

func main() {
	log := setupPrettySlog()

	log.Info(
		"starting rest-api-parser",
	)
	log.Debug("debug messages are enabled")

	storage, err := psql.New()
	if err != nil {
		log.Error("failed to initialize storage", slog.String("error", err.Error()))
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/save", save.New(log, storage))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  time.Duration(4 * time.Second),
		WriteTimeout: time.Duration(4 * time.Second),
		IdleTimeout:  time.Duration(30 * time.Second),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("error", err.Error()))

		return
	}

	log.Info("server stopped")
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
