package RESTHandlers

import (
	"net/http"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (h *Handler) GetObject(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.WithField("requestID", middleware.GetReqID(r.Context()))

	id := chi.URLParam(r, "id")

	obj, err := h.Store.Get(id)
	if err != nil {
		if err == store.KeyNotFound {
			log.WithField("objectID", id).Debug(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", obj.MediaType)
	err = h.SendResponse(w, obj.Value)
	if err != nil {
		log.Info(err)
	}
}
