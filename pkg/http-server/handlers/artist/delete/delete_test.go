package delete_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/leyl1ne/rest-api-parser/pkg/api/response"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/delete"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/artist/delete/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		takeID    string
		respError string
		mockError error
	}{
		{
			name:   "Success",
			takeID: "1",
		},
		{
			name:      "Invalid ID",
			takeID:    "aab",
			respError: "invalid artist ID",
		},
		{
			name:      "Artist Not Found",
			takeID:    "1",
			respError: "artist not found",
			mockError: storage.ErrArtistNotFound,
		},
		{
			name:      "Failed To Delete",
			takeID:    "1",
			respError: "failed to delete artist",
			mockError: errors.New("unexpected error"),
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artistDeleterMock := mocks.NewArtistDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				artistDeleterMock.On("DeleteArtist", 1).Return(tc.mockError).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), artistDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.takeID), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/{id}", handler.ServeHTTP)

			r.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
