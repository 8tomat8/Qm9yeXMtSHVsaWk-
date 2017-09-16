package api

import (
	"time"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/RESTHandlers"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/middlewares"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
)

const defaultTimeout = 60

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // Could be improved to use logrus lib to log errors
	r.Use(middleware.Timeout(defaultTimeout * time.Second))
	r.Use(middleware.RequestID)

	val := valve.New()
	r.Mount("/api", apiHandler(val))

	return r
}

// apiHandler will handle all requests for /api path
func apiHandler(val *valve.Valve) chi.Router {
	r := chi.NewRouter()
	r.Mount("/objects", objectsHandler(val))
	return r
}

func objectsHandler(val *valve.Valve) chi.Router {
	r := chi.NewRouter()

	RESTHandler := RESTHandlers.Handler{Store: store.GetStorage(), Valve: val}

	r.With(middlewares.ContentTypeChecker).Post("/", RESTHandler.CreateObject)
	r.Get("/", RESTHandler.ListObjects)
	r.Get("/{id}", RESTHandler.GetObject)
	r.Delete("/{id}", RESTHandler.DeleteObject)
	return r
}
