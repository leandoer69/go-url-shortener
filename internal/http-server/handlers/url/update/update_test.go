package update

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang-url-shortener/internal/http-server/handlers/url/update/mocks"
	"golang-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"golang-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		oldAlias  string
		newAlias  string
		respError string
		mockError error
	}{
		{
			name:      "empty old_alias",
			oldAlias:  "",
			newAlias:  "new_alias",
			url:       "https://www.youtube.com/",
			respError: "field OldAlias is not valid",
		},
		{
			name:      "empty new_alias",
			oldAlias:  "old_alias",
			newAlias:  "",
			url:       "https://www.youtube.com/",
			respError: "field NewAlias is not valid",
		},
		{
			name:     "success",
			oldAlias: "old_google",
			newAlias: "new_google",
			url:      "https://www.youtube.com/",
		},

		{
			name:      "invalid url",
			oldAlias:  "old_google",
			newAlias:  "new_google",
			url:       "wrong url",
			respError: "field URL is not a valid URL",
		},

		{
			name:      "empty url",
			oldAlias:  "old_google",
			newAlias:  "new_google",
			url:       "",
			respError: "field URL is not valid",
		},

		{
			name:      "UpdateURL Error",
			oldAlias:  "old_alias",
			newAlias:  "new_alias",
			url:       "https://google.com",
			respError: "failed to update url",
			mockError: errors.New("unexpected error"),
		},

		{
			name:      "url not found",
			oldAlias:  "old_google",
			newAlias:  "new_google",
			url:       "https://google.com",
			respError: "url with this alias not found",
			mockError: storage.ErrUrlNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUrlUpdater := updatemock.NewMockURLUpdater(ctrl)

			if tc.mockError != nil || tc.respError == "" {
				mockUrlUpdater.EXPECT().UpdateURL(tc.url, tc.oldAlias, tc.newAlias).Return(tc.mockError).Times(1)
			}

			handler := New(slogdiscard.NewDiscardLogger(), mockUrlUpdater)

			input := fmt.Sprintf(`{"url": "%s", "old_alias": "%s", "new_alias": "%s"}`, tc.url, tc.oldAlias, tc.newAlias)

			req, err := http.NewRequest(http.MethodPut, "/url/", bytes.NewBuffer([]byte(input)))
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
