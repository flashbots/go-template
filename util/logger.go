package util

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const LoggerContextKey contextKey = "logger"

func GetLogger(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	if l, found := ctx.Value(LoggerContextKey).(*zap.Logger); found {
		return l
	}
	fallback.Warn("Failed to get the context-specific logger, resorting to the fallback one",
		zap.Stack("stacktrace"),
	)
	return fallback
}
