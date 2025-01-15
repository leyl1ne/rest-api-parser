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
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/delete"
	"github.com/leyl1ne/rest-api-parser/pkg/http-server/handlers/song/delete/mocks"
	"github.com/leyl1ne/rest-api-parser/pkg/logger/handlers/slogdiscard"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		takeId    string
		respError string
		mockError error
	}{
		{
			name:   "Success",
			takeId: "1",
		},
		{
			name:      "Invalid ID",
			takeId:    "ab",
			respError: "invalid song ID",
		},
		{
			name:      "Song Not Found",
			takeId:    "1",
			respError: "song not found",
			mockError: storage.ErrSongNotFound,
		},
		{
			name:      "Failed To Delete",
			takeId:    "1",
			respError: "failed to delete song",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			songDeleterMock := mocks.NewSongDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				songDeleterMock.On("DeleteSong", 1).Return(tc.mockError).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), songDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.takeId), nil)
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
