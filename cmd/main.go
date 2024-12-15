package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leyl1ne/rest-api-parser/pkg/handlers"
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
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	storage, err := psql.New()
	if err != nil {
		log.Error("failed to initialize storage", err)
	}
	handleRequests(DB)
	psql.CloseConnection(DB)
}
