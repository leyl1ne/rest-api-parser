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
	Artist models.Artist
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ArtistGetter
type ArtistGetter interface {
	GetArtist(id int) (models.Artist, error)
}

func New(log *slog.Logger, artistGetter ArtistGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.artist.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		artistID := chi.URLParam(r, "id")
		if artistID == "" {
			log.Error("failed to get ID")

			render.JSON(w, r, resp.Error("artist ID is required"))

			return
		}

		id, err := strconv.Atoi(artistID)
		if err != nil {
			log.Error("failed to conver ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid artist ID"))

			return
		}

		receivedArtist, err := artistGetter.GetArtist(id)
		if errors.Is(err, storage.ErrArtistNotFound) {
			log.Info("Artist not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("artist not found"))

			return
		}
		if err != nil {
			log.Error("failed to get artist", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get artist"))

			return
		}

		log.Info("artist received", slog.Int("id", id))
		responseOK(w, r, receivedArtist)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, artist models.Artist) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Artist:   artist,
	})
}
