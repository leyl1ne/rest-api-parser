package models

import "time"

type Song struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Artist      string        `json:"artist"`
	Album       string        `json:"album, omitempty"`
	ReleaseYear int           `json:"release_year, omitempty"`
	Genre       string        `json:"genre, omitempty"`
	Duration    time.Duration `json:"duration, omitempty"`
	Lyrics      string        `json:"lyrycs"`
}
