package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shorter/internal/lib/api/response"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/internal/lib/random"
	"url-shorter/internal/storage"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty" validate:"url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}
type URLGetter interface {
	GetAllAliases() ([]string, error)
}

func New(log *slog.Logger, urlSaver URLSaver, urlGetter URLGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.url.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("decode was success", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		aliases, err := urlGetter.GetAllAliases()
		if err != nil {
			slog.String("op", op)
			log.Error("failed to fetch all aliases", sl.Err(err))
			return
		}
		if alias == "" {
			alias = random.NewRandomString(aliasLength, aliases)
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}
		log.Info("url saved", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}

}
