package update

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"golang-url-shortener/internal/lib/api/response"
	"golang-url-shortener/internal/lib/logger/sl"
	"golang-url-shortener/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

type Request struct {
	URL      string `json:"url" validate:"required,url"`
	OldAlias string `json:"old_alias" validate:"required"`
	NewAlias string `json:"new_alias" validate:"required"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate mockgen -source=update.go -destination=mocks/updatemock.go -package=updatemock
type URLUpdater interface {
	UpdateURL(urlToUpdate, oldAlias, newAlias string) error
}

func New(log *slog.Logger, urlUpdater URLUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = urlUpdater.UpdateURL(req.URL, req.OldAlias, req.NewAlias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info(
				"url with this alias not found",
				slog.String("url", req.URL),
				slog.String("old_alias", req.OldAlias))

			render.JSON(w, r, response.Error("url with this alias not found"))

			return
		}

		if err != nil {
			log.Error("failed to update url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to update url"))
			return
		}

		log.Info("url updated", slog.String("url", req.URL), slog.String("alias", req.NewAlias))

		responseOK(w, r, req.NewAlias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
