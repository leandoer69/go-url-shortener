package delete

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang-url-shortener/internal/http-server/handlers/url/delete/mocks"
	"golang-url-shortener/internal/lib/api/response"
	"golang-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"golang-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "correct",
			alias: "youtube",
		},
		{
			name:      "url doesn't exist",
			alias:     "youtube",
			respError: "url not found",
			mockError: storage.ErrUrlNotFound,
		},
		{
			name:      "error with db",
			alias:     "youtube",
			respError: "internal error",
			mockError: errors.New("another error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUrlDeleter := mocks.NewMockURLDeleter(ctrl)

			if tc.mockError != nil || tc.respError == "" {
				mockUrlDeleter.EXPECT().DeleteURL(tc.alias).Return(tc.mockError)
			}

			handler := New(slogdiscard.NewDiscardLogger(), mockUrlDeleter)
			router := chi.NewRouter()
			router.Delete("/url/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp response.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
