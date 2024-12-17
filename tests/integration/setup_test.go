package integration

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang-url-shortener/internal/config"
	"golang-url-shortener/internal/http-server/handlers/redirect"
	"golang-url-shortener/internal/http-server/handlers/url/delete"
	"golang-url-shortener/internal/http-server/handlers/url/save"
	"golang-url-shortener/internal/http-server/handlers/url/update"
	"golang-url-shortener/internal/http-server/middleware/logger"
	"golang-url-shortener/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type UrlShortenerSuite struct {
	suite.Suite
	test       *assert.Assertions
	storage    *sqlite.Storage
	server     *httptest.Server
	httpClient *http.Client
}

func TestUrlShortenerSuite(t *testing.T) {
	suite.Run(t, new(UrlShortenerSuite))
}

func (s *UrlShortenerSuite) SetupTest() {
	s.T().Helper()
	s.test = assert.New(s.T())

	cfg := config.Config{
		Env:         "test",
		StoragePath: "./storage.db",
		HTTPServer: config.HTTPServer{
			Address: "localhost:8080",
		},
	}

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		os.Exit(1)
	}

	router := s.setupRouter(storage)
	s.storage = storage
	s.server = httptest.NewServer(router)
	s.httpClient = &http.Client{}
}

func (s *UrlShortenerSuite) TearDownTest() {
	s.T().Helper()

	s.server.Close()

	s.test.NoError(s.storage.ClearDB())
}

func (s *UrlShortenerSuite) setupRouter(storage *sqlite.Storage) *chi.Mux {
	router := chi.NewRouter()

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	router.Use(middleware.RequestID)
	router.Use(logger.New(nopLogger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Post("/", save.New(nopLogger, storage))
		r.Delete("/{alias}", delete.New(nopLogger, storage))
		r.Put("/", update.New(nopLogger, storage))
	})

	router.Get("/{alias}", redirect.New(nopLogger, storage))

	return router
}
