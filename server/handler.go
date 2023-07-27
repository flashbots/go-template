package server

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
	s.isReadyMx.RLock()
	defer s.isReadyMx.RUnlock()

	if !s.isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDrain(w http.ResponseWriter, r *http.Request) {
	l := logutils.ZapFromRequest(r)

	s.isReadyMx.Lock()
	if !s.isReady {
		s.isReadyMx.Unlock()
		return
	}

	s.isReady = false
	l.Info("Server marked as not ready")

	// Let's not hold onto the lock in our sleep
	s.isReadyMx.Unlock()

	// Give LB enough time to detect us unhealthy
	time.Sleep(s.cfg.DrainDuration)
}

func (s *Server) handleUndrain(w http.ResponseWriter, r *http.Request) {
	l := logutils.ZapFromRequest(r)

	s.isReadyMx.Lock()
	defer s.isReadyMx.Unlock()

	if s.isReady {
		return
	}

	s.isReady = true
	l.Info("Server marked as ready")
}
