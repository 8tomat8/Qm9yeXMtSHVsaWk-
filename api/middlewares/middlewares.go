package middlewares

import (
	"context"
	"net/http"

	"mime"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/go-chi/chi/middleware"
)

const (
	MediaTypeKey   = "MediaTypeKey"
	ContentTypeKey = "Content-Type"
)

func ContentType(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get(ContentTypeKey))

		if err != nil || len(mediaType) == 0 {
			log := logger.Log.WithField("requestID", middleware.GetReqID(r.Context()))
			log.Debug("Received request without media type")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, MediaTypeKey, mediaType)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetMediaType(r *http.Request) string {
	res, ok := r.Context().Value(MediaTypeKey).(string)
	if !ok {
		return ""
	}
	return res
}
