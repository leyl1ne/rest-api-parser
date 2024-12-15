package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("song ID is required")))
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid song ID")))
	}

	queryStmt := `DELETE FROM songs WHERE id = $1;`
	_, err = h.DB.Query(queryStmt, &id)
	if err != nil {
		log.Println("failed to delete song", err)
		render.JSON(w, r, ErrorRenderer(fmt.Errorf("failed to delte song")))
		return
	}

	render.JSON(w, r, "Deleted")
}
