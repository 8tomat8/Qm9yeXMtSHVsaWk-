package RESTHandlers

import (
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/middlewares"
)

// apiHandler will handle all requests for /api path
func APIHandler(storage store.Storage) chi.Router {
	r := chi.NewRouter()
	r.Mount("/objects", objectsHandler(storage))
	return r
}

// objectsHandler will handle all requests for /objects path
func objectsHandler(storage store.Storage) chi.Router {
	r := chi.NewRouter()

	RESTHandler := Handler{Store: storage}

	r.With(middlewares.ContentType).Post("/", RESTHandler.CreateObject)
	r.Get("/{id}", RESTHandler.GetObject)
	r.Delete("/{id}", RESTHandler.DeleteObject)
	return r
}
