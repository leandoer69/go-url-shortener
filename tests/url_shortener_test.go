package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"golang-url-shortener/internal/http-server/handlers/url/save"
	"golang-url-shortener/internal/lib/random"
	"net/url"
	"testing"
)

var (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/url",
	}

	e := httpexpect.Default(t, u.String())
	e.POST("/").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}
