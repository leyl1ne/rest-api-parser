package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/update"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/update/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	cases := []struct {
		name       string
		takeId     string
		updateSong models.Song
		respError  string
		mockError  error
	}{

		{
			name:   "Success",
			takeId: "12",
			updateSong: models.Song{
				Title:       "Кот",
				Artist:      "Бар Хороших Людей",
				Album:       "Брянский шум",
				ReleaseYear: 2024,
				Genre:       "Гринж",
				Lyrics:      "Мяу, Мяу, Мяу",
			},
		},
		{
			name:   "Empty ID",
			takeId: "",
			updateSong: models.Song{
				Title:       "Кот",
				Artist:      "Бар Хороших Людей",
				Album:       "Брянский шум",
				ReleaseYear: 2024,
				Genre:       "Гринж",
				Lyrics:      "Мяу, Мяу, Мяу",
			},
			respError: "song ID is required",
		},
		{
			name:   "ID is not digit",
			takeId: "fd",
			updateSong: models.Song{
				Title:       "Кот",
				Artist:      "Бар Хороших Людей",
				Album:       "Брянский шум",
				ReleaseYear: 2024,
				Genre:       "Гринж",
				Lyrics:      "Мяу, Мяу, Мяу",
			},
			respError: "invalid song ID",
		},
		{
			name:      "Empty request",
			takeId:    "1",
			respError: "empty request",
		},
		{
			name:   "UpdateSong Error",
			takeId: "1",
			updateSong: models.Song{
				Title:       "Кот",
				Artist:      "Бар Хороших Людей",
				Album:       "Брянский шум",
				ReleaseYear: 2024,
				Genre:       "Гринж",
				Lyrics:      "Мяу, Мяу, Мяу",
			},
			respError: "failed to update song",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			songUpdaterMock := mocks.NewSongUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				songUpdaterMock.On("UpdateSong", mock.AnythingOfType("int"), tc.updateSong).
					Return(tc.mockError).
					Once()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), songUpdaterMock)

			inputString := fmt.Sprintf(`{"title": "%s", "artist": "%s", "album": "%s", "release_year": %d, "genre": "%s", "lyrics": "%s"}`,
				tc.updateSong.Title, tc.updateSong.Artist, tc.updateSong.Album, tc.updateSong.ReleaseYear, tc.updateSong.Genre, tc.updateSong.Lyrics)
			input := bytes.NewReader([]byte(inputString))

			if tc.name == "Empty request" {
				input = bytes.NewReader([]byte{})
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s", tc.takeId), input)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/{id}", handler.ServeHTTP)

			r.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}