package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/save"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/save/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		song      models.Song
		respError string
		mockError error
	}{
		{
			name: "Success",
			song: models.Song{
				Title:       "Кот",
				Artist:      "Бар Хороших Людей",
				Album:       "Брянский шум",
				ReleaseYear: 2024,
				Genre:       "Гринж",
				Lyrics:      "Мяу, Мяу, Мяу",
			},
		},
		{
			name: "Empty album, releaseYear, genre",
			song: models.Song{
				Title:  "Кот",
				Artist: "Бар Хороших Людей",
				Lyrics: "Мяу, Мяу, Мяу",
			},
		},
		{
			name: "Empty title",
			song: models.Song{
				Artist: "Бар Хороших Людей",
				Lyrics: "Мяу, Мяу, Мяу",
			},
			respError: "field Title is required field",
		},
		{
			name:      "Empty request",
			respError: "empty request",
		},
		{
			name: "SaveSong Error",
			song: models.Song{
				Title:  "Кот",
				Artist: "Бар Хороших Людей",
				Lyrics: "Мяу, Мяу, Мяу",
			},
			respError: "failed to add song",
			mockError: errors.New("unexpected error"),
		}}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			songSaverMock := mocks.NewSongSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				songSaverMock.On("SaveSong", tc.song).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), songSaverMock)

			inputString := fmt.Sprintf(`{"title": "%s", "artist": "%s", "album": "%s", "release_year": %d, "genre": "%s", "lyrics": "%s"}`,
				tc.song.Title, tc.song.Artist, tc.song.Album, tc.song.ReleaseYear, tc.song.Genre, tc.song.Lyrics)
			input := bytes.NewReader([]byte(inputString))

			if tc.name == "Empty request" {
				input = bytes.NewReader([]byte{})
			}

			req, err := http.NewRequest(http.MethodPost, "/save", input)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
