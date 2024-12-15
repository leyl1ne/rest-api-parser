package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

// TODO: modify the handler by chi render
func (h handler) GetAllSong(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query("SELECT * FROM articles;")
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var songs = make([]models.Song, 0)
	for results.Next() {
		var song models.Song
		err = results.Scan(&song.ID, &song.Title,
			&song.Artist, &song.Album, &song.ReleaseYear,
			&song.Genre, &song.Duration, &song.Lyrics)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}

		songs = append(songs, song)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}
