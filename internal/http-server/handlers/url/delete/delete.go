package delete

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shorter/internal/lib/api/response"
	"url-shorter/internal/lib/logger/sl"
)

type Response struct {
	resp.Response
}
type Request struct {
	Alias string `json:"alias" validate:"url,required"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.delete.url.New"

		log = log.With(
			slog.String("op", op),
			slog.String("try to delete", r.URL.Path),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		err = urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to delete url by alias %v", err)))

			return
		}

		log.Info("urls deleted")
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
