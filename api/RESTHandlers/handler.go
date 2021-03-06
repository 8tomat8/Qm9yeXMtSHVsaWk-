package RESTHandlers

import (
	"net/http"

	"errors"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
)

type Handler struct {
	Store store.Storage
}

func (Handler) SendResponse(w http.ResponseWriter, payload []byte) error {
	n, err := w.Write(payload)
	if err != nil {
		return err
	} else if n != len(payload) {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("not all data was written to response")
	}
	return nil
}
