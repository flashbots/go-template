package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/flashbots/go-template/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/atomic"
)

type Server struct {
	cfg     *Config
	id      uuid.UUID
	isReady atomic.Bool
	log     *slog.Logger
	srv     *http.Server
}

func New(cfg *Config) *Server {
	id := uuid.Must(uuid.NewRandom())

	s := &Server{
		cfg:     cfg,
		id:      id,
		isReady: atomic.Bool{},
		log:     cfg.Log.With("serverID", id.String()),
		srv:     nil,
	}
	s.isReady.Swap(true)

	mux := chi.NewRouter()
	mux.With(s.httpLogger).Get("/api", s.handleAPI) // Never serve at `/` (root) path
	mux.With(s.httpLogger).Get("/livez", s.handleLivenessCheck)
	mux.With(s.httpLogger).Get("/readyz", s.handleReadinessCheck)
	mux.With(s.httpLogger).Get("/drain", s.handleDrain)
	mux.With(s.httpLogger).Get("/undrain", s.handleUndrain)

	s.srv = &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return s
}

func (s *Server) httpLogger(next http.Handler) http.Handler {
	// return httplogger.LoggingMiddlewareZap(s.log, next)
	return next // TODO: slog logging middleware
}

func (s *Server) RunInBackground() {
	// metrics
	if s.cfg.MetricsAddr != "" {
		s.log.With("metricsAddress", s.cfg.MetricsAddr).Info("Starting metrics server")
		go func() {
			if err := metrics.ListenAndServe(
				"github.com/flashbots/go-template",
				s.cfg.MetricsAddr,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("HTTP server failed", "err", err)
			}
		}()
	}

	// api
	{
		s.log.Info("Starting HTTP server",
			slog.String("listenAddress", s.cfg.ListenAddr),
			slog.String("version", s.cfg.Version),
		)

		go func() {
			if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("HTTP server failed", "err", err)
			}
		}()
	}
}

func (s *Server) Shutdown() {
	// api
	{
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulShutdownDuration)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			s.log.Error("Graceful HTTP server shutdown failed", "err", err)
		} else {
			s.log.Info("HTTP server gracefully stopped")
		}
	}

	// metrics
	if len(s.cfg.MetricsAddr) != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulShutdownDuration)
		defer cancel()

		if err := metrics.Shutdown(ctx); err != nil {
			s.log.Error("Graceful metrics server shutdown failed", "err", err)
		} else {
			s.log.Info("Metrics server gracefully stopped")
		}
	}
}
