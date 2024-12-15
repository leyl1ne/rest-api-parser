package mocks

import (
	"github.com/leyl1ne/rest-api-parser/pkg/models"
)

var Songs = []models.Song{
	{
		ID:          1,
		Title:       "Кот",
		Artist:      "Бар хороших людей",
		ReleaseYear: 2024,
		Genre:       "Гринж",
		Duration:    150,
		Lyrics: `Мяу, Мяу, Мяу, Мяу
		Мяу, Мяу, Мяу, Мяу
		Мяу, Мяу, Мяу, Мяу
		Мяу, Мяу, Мяу, Мяу
		Мяу, Мяу, Мяу, Мяу
		Мяу, Мяу, Мяу, Мяу`,
	},
	{
		ID:     2,
		Title:  "Salut",
		Artist: "Stromae",
		Lyrics: `Salut, Salut, Salut, Salut
		Salut, Salut, Salut, Salut
		Salut, Salut, Salut, Salut
		Salut, Salut, Salut, Salut
		Salut, Salut, Salut, Salut`,
	},
}
