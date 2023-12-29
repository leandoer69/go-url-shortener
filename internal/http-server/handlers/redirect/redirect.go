package redirect

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang-url-shortener/internal/lib/logger/sl"
	"golang-url-shortener/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate mockgen -source=redirect.go -destination=mocks/redirectmock.go -package=mocks
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, "invalid request")
			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", sl.Err(err))
			render.JSON(w, r, "url not found")
			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, "internal error")
			return
		}
		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
