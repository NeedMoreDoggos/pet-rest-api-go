package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl"
	resp "github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl/api/response"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGiver interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGiver URLGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGiver.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", sl.Err(err))
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url got and redirected", slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
