package get_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/get"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/get/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name       string
		takeId     string
		returnSong models.Song
		respError  string
		mockError  error
	}{
		{
			name:   "Success",
			takeId: "1",
			returnSong: models.Song{
				Title:  "Slim Shady",
				Artist: "Eminem",
				Lyrics: "Real slim shady",
			},
		},
		{
			name:      "Invalid ID",
			takeId:    "ab",
			respError: "invalid song ID",
		},
		{
			name:      "Not Found Song",
			takeId:    "1",
			respError: "song not found",
			mockError: storage.ErrSongNotFound,
		},
		{
			name:      "Failed to update",
			takeId:    "1",
			respError: "failed to get song",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			songGetterMock := mocks.NewSongGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				songGetterMock.On("GetSong", 1).
					Return(tc.returnSong, tc.mockError).
					Once()
			}

			handler := get.New(slogdiscard.NewDiscardLogger(), songGetterMock)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.takeId), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/{id}", handler.ServeHTTP)

			r.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp get.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
