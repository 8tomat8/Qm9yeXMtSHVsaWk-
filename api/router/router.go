package router

import (
	"time"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/RESTHandlers"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultTimeout = 60

func NewRouter(storage store.Storage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // Could be improved to use custom lib to log errors
	r.Use(middleware.Timeout(defaultTimeout * time.Second))
	r.Use(middleware.RequestID)

	r.Mount("/api", RESTHandlers.APIHandler(storage))

	return r
}
