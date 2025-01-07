package getAll_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/getAll"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/getAll/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestGetAllHandler(t *testing.T) {
	cases := []struct {
		name        string
		returnSongs []models.Song
		respError   string
		mockError   error
	}{
		{
			name: "Success",
			returnSongs: []models.Song{
				{
					Title:  "Кот",
					Artist: "Бар Хороших Людей",
					Lyrics: "Мяу, Мяу, Мяу",
				},
				{
					Title:  "Кот",
					Artist: "Бар Хороших Людей",
					Lyrics: "Мяу, Мяу, Мяу",
				},
				{
					Title:  "Кот",
					Artist: "Бар Хороших Людей",
					Lyrics: "Мяу, Мяу, Мяу",
				},
			},
		},
		{
			name:      "Failed To Get",
			respError: "failed to get all song",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			songGetterAllMock := mocks.NewSongGetterAll(t)

			if tc.respError == "" || tc.mockError != nil {
				songGetterAllMock.On("GetAllSong").Return(tc.returnSongs, tc.mockError).Once()
			}

			handler := getAll.New(slogdiscard.NewDiscardLogger(), songGetterAllMock)

			req, err := http.NewRequest(http.MethodGet, "/getAll", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp getAll.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			require.Equal(t, tc.returnSongs, resp.Songs)
		})
	}
}
