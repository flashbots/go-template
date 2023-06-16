package server

import (
	"context"
	"net/http"

	"github.com/flashbots/go-template/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Server) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Setup the logger
		var l *zap.Logger
		if logger, found := ctx.Value(util.LoggerContextKey).(*zap.Logger); found {
			l = logger
		} else {
			l = s.log
		}

		// Add common fields
		l = l.With(
			zap.String("httpRequestID", uuid.New().String()),
		)

		// Inject logger into context
		ctx = context.WithValue(ctx, util.LoggerContextKey, l)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
