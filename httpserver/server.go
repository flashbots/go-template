package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	victoriaMetrics "github.com/VictoriaMetrics/metrics"
	"github.com/flashbots/go-template/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"go.uber.org/atomic"
)

type HTTPServerConfig struct {
	ListenAddr  string
	MetricsAddr string
	EnablePprof bool
	Log         *httplog.Logger

	DrainDuration            time.Duration
	GracefulShutdownDuration time.Duration
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
}

type Server struct {
	cfg     *HTTPServerConfig
	isReady atomic.Bool
	log     *httplog.Logger

	srv        *http.Server
	metricsSrv *http.Server
}

func New(cfg *HTTPServerConfig) (srv *Server, err error) {
	srv = &Server{
		cfg: cfg,
		log: cfg.Log,
		srv: nil,
	}

	if cfg.MetricsAddr != "" {
		srv.metricsSrv = &http.Server{
			Addr:         cfg.MetricsAddr,
			Handler:      srv.getMetricsRouter(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}
	}

	srv.srv = &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      srv.getRouter(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	srv.isReady.Swap(true)

	return srv, nil
}

func (srv *Server) getRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(httplog.RequestLogger(srv.log))
	mux.Use(middleware.Recoverer)
	mux.Use(metrics.Middleware)

	mux.Get("/api", srv.handleAPI) // Never serve at `/` (root) path
	mux.Get("/livez", srv.handleLivenessCheck)
	mux.Get("/readyz", srv.handleReadinessCheck)
	mux.Get("/drain", srv.handleDrain)
	mux.Get("/undrain", srv.handleUndrain)

	if srv.cfg.EnablePprof {
		srv.log.Info("pprof API enabled")
		mux.Mount("/debug", middleware.Profiler())
	}
	return mux
}

func (srv *Server) getMetricsRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		victoriaMetrics.WritePrometheus(w, true)
	})
	return mux
}

func (srv *Server) RunInBackground() {
	// metrics
	if srv.cfg.MetricsAddr != "" {
		go func() {
			srv.log.With("metricsAddress", srv.cfg.MetricsAddr).Info("Starting metrics server")
			err := srv.metricsSrv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				srv.log.Error("HTTP server failed", "err", err)
			}
		}()
	}

	// api
	go func() {
		srv.log.Info("Starting HTTP server", "listenAddress", srv.cfg.ListenAddr)
		if err := srv.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srv.log.Error("HTTP server failed", "err", err)
		}
	}()
}

func (srv *Server) Shutdown() {
	// api
	ctx, cancel := context.WithTimeout(context.Background(), srv.cfg.GracefulShutdownDuration)
	defer cancel()
	if err := srv.srv.Shutdown(ctx); err != nil {
		srv.log.Error("Graceful HTTP server shutdown failed", "err", err)
	} else {
		srv.log.Info("HTTP server gracefully stopped")
	}

	// metrics
	if len(srv.cfg.MetricsAddr) != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), srv.cfg.GracefulShutdownDuration)
		defer cancel()

		if err := srv.metricsSrv.Shutdown(ctx); err != nil {
			srv.log.Error("Graceful metrics server shutdown failed", "err", err)
		} else {
			srv.log.Info("Metrics server gracefully stopped")
		}
	}
}
