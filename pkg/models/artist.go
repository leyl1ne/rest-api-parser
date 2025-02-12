package models

type Artist struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name" validate:"required"`
	Genre   string `json:"genre,omitempty"`
	Country string `json:"country,omitempty"`
}
