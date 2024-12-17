package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongSaver
type SongSaver interface {
	SaveSong(song models.Song) (int64, error)
}

func New(log *slog.Logger, songSaver SongSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var song models.Song
		err := render.DecodeJSON(r.Body, &song)

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

		if err := validator.New().Struct(song); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := songSaver.SaveSong(song)
		if errors.Is(err, storage.ErrSongExists) {
			log.Info("Song already exists", slog.Int("id", song.ID))

			render.JSON(w, r, resp.Error("song already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add song")

			render.JSON(w, r, resp.Error("failed to add song"))

			return
		}

		log.Info("song added", slog.Int64("id", id))

		render.JSON(w, r, resp.OK())
	}
}
