package RESTHandlers

import (
	"net/http"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (h *Handler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log := logger.Log.WithField("requestID", middleware.GetReqID(r.Context()))

	err := h.Store.Delete(id)
	if err != nil {
		if err == store.KeyNotFound {
			log.WithField("id", id).Debug(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
