package RESTHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/go-chi/chi/middleware"
)

// TODO: REMOVE IT!!!! FOR DEV NEEDS ONLY!!!
func (h *Handler) ListObjects(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.WithField("requestID", middleware.GetReqID(r.Context()))

	objs := h.Store.GetAll()
	jsonedObjs, err := json.Marshal(objs)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(jsonedObjs)
	if err != nil {
		log.Error(err)
		return
	} else if n != len(jsonedObjs) {
		log.Error("Not all data was successfully written to response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
