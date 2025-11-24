package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shorter/internal/lib/api/response"
	"url-shorter/internal/storage"
)

type Response struct {
	resp.Response
	URL string `json:"url"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect.url.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := "https://" + chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}
		resUrl, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrorURLNotFound) {

			render.JSON(w, r, resp.Error("url not found"))

			return
		}
		if err != nil {
			log.Info("inernal error")
			render.JSON(w, r, resp.Error("inernal error"))

			return
		}
		log.Info("url found", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)
	}

}
