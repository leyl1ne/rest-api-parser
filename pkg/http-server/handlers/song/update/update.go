package update

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongUpdater
type SongUpdater interface {
	UpdateSong(id int, updatedSong models.Song) error
}

func New(log *slog.Logger, songUpdater SongUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.update.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songID := chi.URLParam(r, "id")
		if songID == "" {
			log.Error("filed to get ID")

			render.JSON(w, r, resp.Error("song ID is required"))

			return
		}

		id, err := strconv.Atoi(songID)
		if err != nil {
			log.Error("failed to convert ID", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("invalid song ID"))
		}

		var updatedSong models.Song
		err = render.DecodeJSON(r.Body, &updatedSong)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded")

		err = songUpdater.UpdateSong(id, updatedSong)
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Info("Song not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("song not found"))

			return
		}
		if err != nil {
			log.Error("failed to update song")

			render.JSON(w, r, resp.Error("failed to update song"))

			return
		}

		log.Info("song updated", slog.Int("id", id))
		render.JSON(w, r, resp.OK())
	}
}
