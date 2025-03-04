package getall_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	getall "github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/getAll"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/getAll/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestGetAllGandler(t *testing.T) {
	cases := []struct {
		name          string
		returnArtists []models.Artist
		respError     string
		mockError     error
	}{
		{
			name: "Success",
			returnArtists: []models.Artist{
				{
					Name:    "Бар Хороших Людей",
					Genre:   "Кринж",
					Country: "Russia",
				},
				{
					Name:    "Бар Хороших Людей",
					Genre:   "Кринж",
					Country: "Russia",
				},
				{
					Name:    "Бар Хороших Людей",
					Genre:   "Кринж",
					Country: "Russia",
				},
			},
		},
		{
			name:      "Failed to get",
			respError: "failed to get all artists",
			mockError: errors.New("unexpected error"),
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artistGetterAllMock := mocks.NewArtistGetterAll(t)

			if tc.respError == "" || tc.mockError != nil {
				artistGetterAllMock.On("GetAllArtist").Return(tc.returnArtists, tc.mockError).Once()
			}

			handler := getall.New(slogdiscard.NewDiscardLogger(), artistGetterAllMock)

			req, err := http.NewRequest(http.MethodGet, "/getAll", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp getall.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			require.Equal(t, tc.returnArtists, resp.Artist)
		})
	}
}
