package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl/api/response"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
}

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
		if errors.Is(err, storage.ErrURLNotFound) {
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
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
