package get_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/get"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/get/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name         string
		takeID       string
		returnArtist models.Artist
		respError    string
		mockError    error
	}{
		{
			name:   "Success",
			takeID: "1",
			returnArtist: models.Artist{
				Name:    "Бар Хороших Людей",
				Genre:   "Кринж",
				Country: "Russia",
			},
		},
		{
			name:      "Invalid ID",
			takeID:    "ab",
			respError: "invalid artist ID",
		},
		{
			name:      "Not Found Song",
			takeID:    "1",
			respError: "artist not found",
			mockError: storage.ErrArtistNotFound,
		},
		{
			name:      "Failed to get",
			takeID:    "1",
			respError: "failed to get artist",
			mockError: errors.New("unexpected error"),
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artistGetterMock := mocks.NewArtistGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				artistGetterMock.On("GetArtist", 1).Return(tc.returnArtist, tc.mockError).Once()
			}

			handler := get.New(slogdiscard.NewDiscardLogger(), artistGetterMock)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.takeID), nil)
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
