package redirect

import (
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang-url-shortener/internal/http-server/handlers/redirect/mocks"
	"golang-url-shortener/internal/lib/api"
	"golang-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "correct",
			alias: "youtube",
			url:   "https://www.youtube.com/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockUrlDeleter := mocks.NewMockURLGetter(ctrl)

			if tc.mockError != nil || tc.respError == "" {
				mockUrlDeleter.EXPECT().GetURL(tc.alias).Return(tc.url, tc.mockError).Times(1)
			}

			r := chi.NewRouter()
			r.Get("/{alias}", New(slogdiscard.NewDiscardLogger(), mockUrlDeleter))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// check if we got redirected
			require.Equal(t, redirectedToUrl, tc.url)
		})
	}
}
