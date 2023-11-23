package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/flashbots/go-template/metrics"
	"github.com/flashbots/go-utils/httplogger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Server struct {
	cfg     *Config
	id      uuid.UUID
	isReady atomic.Bool
	log     *zap.Logger
	srv     *http.Server
}

func New(cfg *Config) *Server {
	id := uuid.Must(uuid.NewRandom())
	s := &Server{
		cfg:     cfg,
		id:      id,
		isReady: atomic.Bool{},
		log:     cfg.Log.With(zap.String("serverID", id.String())),
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
	return httplogger.LoggingMiddlewareZap(s.log, next)
}

func (s *Server) RunInBackground() {
	// metrics
	if len(s.cfg.MetricsAddr) != 0 {
		s.log.Info("Starting metrics server",
			zap.String("metricsAddress", s.cfg.MetricsAddr),
		)
		metrics.Init("github.com/flashbots/go-template", s.cfg.MetricsAddr)
		go func() {
			if err := metrics.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("HTTP server failed", zap.Error(err))
			}
		}()
	}

	// api
	{
		s.log.Info("Starting HTTP server",
			zap.String("listenAddress", s.cfg.ListenAddr),
			zap.String("version", s.cfg.Version),
		)

		go func() {
			if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("HTTP server failed", zap.Error(err))
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
			s.log.Error("Graceful HTTP server shutdown failed", zap.Error(err))
		} else {
			s.log.Info("HTTP server gracefully stopped")
		}
	}

	// metrics
	if len(s.cfg.MetricsAddr) != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulShutdownDuration)
		defer cancel()

		if err := metrics.Shutdown(ctx); err != nil {
			s.log.Error("Graceful metrics server shutdown failed", zap.Error(err))
		} else {
			s.log.Info("Metrics server gracefully stopped")
		}
	}
}
