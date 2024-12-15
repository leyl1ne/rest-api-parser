package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers"
	mwLogger "github.com/leyl1ne/rest-api-parser/pkg/http-server/middleware/logger"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogpretty"
	"github.com/leyl1ne/rest-api-parser/pkg/storage/psql"
)

func handleRequests(DB *sql.DB) {
	h := handlers.New(DB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", h.GetAllSong)
	r.Get("/{id}", h.GetSong)
	r.Post("/add", h.AddSong)
	r.Put("/add/{id}", h.UpdateSong)
	r.Delete("/delete/{id}", h.DeleteSong)
}

func main() {
	log := setupPrettySlog()

	log.Info(
		"starting rest-api-parser",
	)
	log.Debug("debug messages are enabled")

	storage, err := psql.New()
	if err != nil {
		log.Error("failed to initialize storage", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	psql.CloseConnection()
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
