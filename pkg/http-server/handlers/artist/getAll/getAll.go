package getall

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

type Response struct {
	resp.Response
	Artist []models.Artist
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ArtistGetterAll
type ArtistGetterAll interface {
	GetAllArtist() ([]models.Artist, error)
}

func New(log *slog.Logger, artistGetterAll ArtistGetterAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.artist.getAll.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		artists, err := artistGetterAll.GetAllArtist()
		if err != nil {
			log.Error("failed to get all artist", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get all artists"))

			return
		}

		log.Info("artists received")
		responseOK(w, r, artists)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, artists []models.Artist) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Artist:   artists,
	})
}
