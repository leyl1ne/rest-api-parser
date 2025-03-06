package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ArtistDeleter
type ArtistDeleter interface {
	DeleteArtist(id int) error
}

func New(log *slog.Logger, artistDeleter ArtistDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.song.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		artistID := chi.URLParam(r, "id")
		if artistID == "" {
			log.Error("failed to get ID")

			render.JSON(w, r, resp.Error("invalid aritst ID"))

			return
		}

		id, err := strconv.Atoi(artistID)
		if err != nil {
			log.Error("failed to conver ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid artist ID"))

			return
		}

		err = artistDeleter.DeleteArtist(id)
		if errors.Is(err, storage.ErrArtistNotFound) {
			log.Info("artist not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("artist not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete artist", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to delete artist"))

			return
		}

		log.Info("aritst delete", slog.Int("id", id))

		render.JSON(w, r, resp.OK())
	}
}
