package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

func (h handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("song ID is required")))
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid song ID")))
	}

	var updatedSong models.Song
	err = render.DecodeJSON(r.Body, &updatedSong)
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

	queryStmt := `UPDATE songs SET 
	title = $2, artist = $3, album = $4, 
	release_year = $5, genre = $6, duration = $7, lyrics = $8 
	WHERE id = $1;`
	err = h.DB.QueryRow(queryStmt, &id, &updatedSong.Title,
		&updatedSong.Artist, &updatedSong.Album, &updatedSong.ReleaseYear,
		&updatedSong.Genre, &updatedSong.Duration, &updatedSong.Lyrics).Scan()
	if err != nil {
		log.Println("failed to update song", err)

		render.JSON(w, r, ErrorRenderer(fmt.Errorf("failed to update song")))

		return
	}

	render.JSON(w, r, "Updated")
}
