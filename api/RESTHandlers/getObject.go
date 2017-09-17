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

	id := chi.URLParamFromCtx(r.Context(), "id")

	obj, err := h.Store.Get(r.Context(), id)
	if err != nil {
		switch err {
		case store.KeyNotFound:
			log.WithField("objectID", id).Debug(err)
			w.WriteHeader(http.StatusNotFound)
		default:
			log.Error(err)
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
