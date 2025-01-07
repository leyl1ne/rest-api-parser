package getAll

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

type Response struct {
	resp.Response
	Songs []models.Song
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SongGetterAll
type SongGetterAll interface {
	GetAllSong() ([]models.Song, error)
}

func New(log *slog.Logger, songGetterAll SongGetterAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.getAll.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songs, err := songGetterAll.GetAllSong()
		if err != nil {

			log.Error("failed to get all song", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get all song"))

			return
		}

		log.Info("songs received")
		responseOK(w, r, songs)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, songs []models.Song) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Songs:    songs,
	})
}
