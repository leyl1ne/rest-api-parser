package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

func (h handler) GetSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("song ID is required")))
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid song ID")))
	}

	queryStmt := `SELECT * FROM songs WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		render.Render(w, r, ErrorRenderer(fmt.Errorf("failed to execute query: %w", err)))
		return
	}

	var song models.Song
	for results.Next() {
		err = results.Scan(&song.ID, &song.Title,
			&song.Artist, &song.Album, &song.ReleaseYear,
			&song.Genre, &song.Duration, &song.Lyrics)
		if err != nil {
			log.Println("failed to scan", err)
			render.Render(w, r, ErrorRenderer(err))
			return
		}
	}

	render.JSON(w, r, song)

}
