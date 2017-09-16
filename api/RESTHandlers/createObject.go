package RESTHandlers

import (
	"net/http"

	"io"

	"strconv"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/middlewares"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi/middleware"
)

func (h *Handler) CreateObject(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.WithField("requestID", middleware.GetReqID(r.Context()))

	contentLen, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Debug("Content-Length could not be parsed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if contentLen > store.MaxValueSize {
		log.Debug("Content-Length more than allowed")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	// This solution could be improved by creating some pool of slices if we can predict loads
	// It will allow us to avoid extra allocations
	value := make([]byte, contentLen)
	_, err = io.ReadFull(r.Body, value)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Warn(err)
		}
	}()

	// Creation process
	key, err := h.Store.Create(r.Context(), value, middlewares.GetMediaType(r))
	if err != nil {
		switch err {
		case store.KeyAlreadyExist:
			log.Debug(err)
			w.WriteHeader(http.StatusConflict)
		case store.ValueSizeIsExceeded:
			log.Warn(err)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		default:
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = h.SendResponse(w, []byte(key))
	if err != nil {
		log.Warn(err)
	}
}
