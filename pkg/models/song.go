package models

type Song struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title" validate:"required"`
	Artist      string `json:"artist" validate:"required"`
	Album       string `json:"album,omitempty"`
	ReleaseYear int    `json:"release_year,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Lyrics      string `json:"lyrics" validate:"required"`
}
