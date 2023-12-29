package delete

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang-url-shortener/internal/lib/api/response"
	"golang-url-shortener/internal/lib/logger/sl"
	"golang-url-shortener/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate mockgen -source=delete.go -destination=mocks/deletemock.go -package=mocks
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", sl.Err(err))
			render.JSON(w, r, response.Error("url not found"))
			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, response.OK())
	}
}
