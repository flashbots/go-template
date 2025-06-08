package httpserver

import (
	"net/http"
	"time"
)

func (srv *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if !srv.isReady.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleDrain(w http.ResponseWriter, r *http.Request) {
	if wasReady := srv.isReady.Swap(false); !wasReady {
		return
	}
	// l := logutils.ZapFromRequest(r)
	srv.log.Info("Server marked as not ready")
	time.Sleep(srv.cfg.DrainDuration) // Give LB enough time to detect us not ready
}

func (srv *Server) handleUndrain(w http.ResponseWriter, r *http.Request) {
	if wasReady := srv.isReady.Swap(true); wasReady {
		return
	}
	// l := logutils.ZapFromRequest(r)
	srv.log.Info("Server marked as ready")
}
