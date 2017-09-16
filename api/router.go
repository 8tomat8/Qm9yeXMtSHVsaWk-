package api

import (
	"time"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/RESTHandlers"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/middlewares"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultTimeout = 60

func NewRouter(storage store.Storage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // Could be improved to use logrus lib to log errors
	r.Use(middleware.Timeout(defaultTimeout * time.Second))
	r.Use(middleware.RequestID)

	r.Mount("/api", apiHandler(storage))

	return r
}

// apiHandler will handle all requests for /api path
func apiHandler(storage store.Storage) chi.Router {
	r := chi.NewRouter()
	r.Mount("/objects", objectsHandler(storage))
	return r
}

func objectsHandler(storage store.Storage) chi.Router {
	r := chi.NewRouter()

	RESTHandler := RESTHandlers.Handler{Store: storage}

	r.With(middlewares.ContentType).Post("/", RESTHandler.CreateObject)
	r.Get("/{id}", RESTHandler.GetObject)
	r.Delete("/{id}", RESTHandler.DeleteObject)
	return r
}
