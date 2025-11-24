package get

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shorter/internal/lib/api/response"
	"url-shorter/internal/lib/logger/sl"
)

type Response struct {
	resp.Response
	Aliases []string `json:"aliases"`
}

type URLGetter interface {
	GetAllAliases() ([]string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.get.url.New"
		log = log.With(
			slog.String("op", op),
			slog.String("try to get alo aliases", r.URL.Path),
		)

		aliases, err := urlGetter.GetAllAliases()
		if err != nil {
			slog.String("op", op)
			log.Error("failed to get all aliases", sl.Err(err))
			return
		}
		log.Info("urls sanded")
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Aliases:  aliases,
		})
	}

}
