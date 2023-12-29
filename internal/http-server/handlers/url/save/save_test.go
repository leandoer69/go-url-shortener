package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang-url-shortener/internal/http-server/handlers/url/save/mocks"
	"golang-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"golang-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveURL(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "empty alias",
			alias: "",
			url:   "https://www.youtube.com/",
		},

		{
			name:  "success",
			alias: "google",
			url:   "https://www.youtube.com/",
		},

		{
			name:      "invalid url",
			alias:     "google",
			url:       "wrong url",
			respError: "field URL is not a valid URL",
		},

		{
			name:      "empty url",
			alias:     "google",
			url:       "",
			respError: "field URL is not valid",
		},

		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},

		{
			name:      "url exists",
			alias:     "google",
			url:       "https://google.com",
			respError: "url already exists",
			mockError: storage.ErrUrlExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUrlSaver := mocks.NewMockURLSaver(ctrl)

			if tc.mockError != nil || tc.respError == "" {
				mockUrlSaver.EXPECT().SaveURL(tc.url, gomock.Any()).Return(int64(1),
					tc.mockError).Times(1)
			}

			handler := New(slogdiscard.NewDiscardLogger(), mockUrlSaver)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url/", bytes.NewBuffer([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
