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
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/save"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/save/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		artist    models.Artist
		respError string
		mockError error
	}{
		{
			name: "Success",
			artist: models.Artist{
				Name:    "Бар Хороших Людей",
				Genre:   "Кринж",
				Country: "Russia",
			},
		},
		{
			name:      "Empty request",
			respError: "empty request",
		},
		{
			name: "Empty other filds, except Name",
			artist: models.Artist{
				Name: "Бар Хороших Людей",
			},
		},
		{
			name: "Empty Name field",
			artist: models.Artist{
				Genre:   "Кринж",
				Country: "Russia",
			},
			respError: "field Name is required field",
		},
		{
			name: "Artist already exists",
			artist: models.Artist{
				Name:    "Бар Хороших Людей",
				Genre:   "Кринж",
				Country: "Russia",
			},
			respError: "artist already exists",
			mockError: storage.ErrArtistExists,
		},
		{
			name: "SaveArtist Error",
			artist: models.Artist{
				Name:    "Бар Хороших Людей",
				Genre:   "Кринж",
				Country: "Russia",
			},
			respError: "failed to add artist",
			mockError: errors.New("unexpected error"),
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artistSaverMock := mocks.NewArtistSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				artistSaverMock.On("SaveArtist", tc.artist).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), artistSaverMock)

			inputString := fmt.Sprintf(`{"name": "%s", "genre": "%s", "country": "%s"}`,
				tc.artist.Name, tc.artist.Genre, tc.artist.Country)
			input := bytes.NewReader([]byte(inputString))

			if tc.name == "Empty request" {
				input = bytes.NewReader([]byte{})
			}

			req, err := http.NewRequest(http.MethodPost, "/save", input)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
