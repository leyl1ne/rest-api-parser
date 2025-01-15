package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongDeleter
type SongDeleter interface {
	DeleteSong(id int) error
}

func New(log *slog.Logger, songDeleter SongDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songID := chi.URLParam(r, "id")
		if songID == "" {
			log.Error("failed to get ID")

			render.JSON(w, r, resp.Error("song ID is required"))

			return
		}

		id, err := strconv.Atoi(songID)
		if err != nil {
			log.Error("failed to convert ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid song ID"))

			return
		}

		err = songDeleter.DeleteSong(id)
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Info("song not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("song not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete song", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to delete song"))

			return
		}

		log.Info("song delete", slog.Int("id", id))

		render.JSON(w, r, resp.OK())
	}
}
