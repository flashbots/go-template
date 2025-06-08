package metrics

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startAt := time.Now()
		defer func() {
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			recordRequestDuration(routePattern, time.Since(startAt).Milliseconds())
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
