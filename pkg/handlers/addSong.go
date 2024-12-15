package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

func (h handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := render.DecodeJSON(r.Body, &song)

	if errors.Is(err, io.EOF) {
		log.Println("request body is empty")

		render.JSON(w, r, ErrorRenderer(fmt.Errorf("empty request")))

		return
	}
	if err != nil {
		log.Println("failed to decode request body", err)

		render.JSON(w, r, ErrorRenderer(fmt.Errorf("failed to decode request body")))

		return
	}

	queryStmt := `INSERT INTO songs (title,artist,album,release_year,genre,duration,lyrics)
			VALUES ($2,$3,$4,$5,$6,$7,$8) RETURNING id;`
	err = h.DB.QueryRow(
		queryStmt, &song.Title,
		&song.Artist, &song.Album, &song.ReleaseYear,
		&song.Genre, &song.Duration, &song.Lyrics).Scan(&song.ID)
	if err != nil {
		log.Println("failed to add url", err)

		render.JSON(w, r, ErrorRenderer(fmt.Errorf("failed to add url")))
	}

	render.JSON(w, r, "Created")
}
