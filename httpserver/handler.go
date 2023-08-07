package httpserver

import (
	"net/http"
	"time"

	"github.com/flashbots/go-utils/logutils"
)

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if !s.isReady.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDrain(w http.ResponseWriter, r *http.Request) {
	if wasReady := s.isReady.Swap(false); !wasReady {
		return
	}
	l := logutils.ZapFromRequest(r)
	l.Info("Server marked as not ready")
	time.Sleep(s.cfg.DrainDuration) // Give LB enough time to detect us not ready
}

func (s *Server) handleUndrain(w http.ResponseWriter, r *http.Request) {
	if wasReady := s.isReady.Swap(true); wasReady {
		return
	}
	l := logutils.ZapFromRequest(r)
	l.Info("Server marked as ready")
}
