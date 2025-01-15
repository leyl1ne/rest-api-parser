package get

import (
	"errors"
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

type Response struct {
	resp.Response
	Song models.Song
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongGetter
type SongGetter interface {
	GetSong(id int) (models.Song, error)
}

func New(log *slog.Logger, songGetter SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get.New"

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

		receivedSong, err := songGetter.GetSong(id)
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Info("Song not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("song not found"))

			return
		}
		if err != nil {
			log.Error("failed to get song", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get song"))

			return
		}

		log.Info("song received", slog.Int("id", id))
		responseOK(w, r, receivedSong)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, song models.Song) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Song:     song,
	})
}
