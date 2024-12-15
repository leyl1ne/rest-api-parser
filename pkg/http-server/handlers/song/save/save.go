package save

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/http-server/errors"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

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

			render.JSON(w, r, resp.Error(fmt.Errorf("empty request")))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error(fmt.Errorf("failed to decode request body")))

			return
		}

		log.Info("request body decoded")

		id, err := songSaver.SaveSong(song)
		if errors.Is(err, storage.ErrSongExists) {
			log.Info("Song already exists", slog.Int("id", song.ID))

			render.JSON(w, r, resp.Error(fmt.Errorf("song already exists")))

			return
		}
		if err != nil {
			log.Error("failed to add song")

			render.JSON(w, r, resp.Error(fmt.Errorf("filed to add song")))

			return
		}

		log.Info("song added", slog.Int64("id", id))

		render.JSON(w, r, "Song added")
	}
}
