package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ArtistSaver
type ArtistSaver interface {
	SaveArtist(artist models.Artist) (int64, error)
}

func New(log *slog.Logger, artistSaver ArtistSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.artist.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var artist models.Artist
		err := render.DecodeJSON(r.Body, &artist)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed  to decode request body"))

			return
		}

		log.Info("request body decoded")

		if err := validator.New().Struct(artist); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := artistSaver.SaveArtist(artist)
		if errors.Is(err, storage.ErrArtistExists) {
			log.Info("Artist already exists", slog.Int("id", artist.ID))

			render.JSON(w, r, resp.Error("artist already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add artist", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to add artist"))

			return
		}

		log.Info("artist added", slog.Int64("id", id))

		render.JSON(w, r, resp.OK())
	}
}
