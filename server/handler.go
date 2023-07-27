package server

import (
	"net/http"
	"time"

	"github.com/flashbots/go-utils/logutils"
)

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	s.isHealthyMx.RLock()
	defer s.isHealthyMx.RUnlock()

	if !s.isHealthy {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// TODO: Add body here if needed
}

func (s *Server) handleDrain(w http.ResponseWriter, r *http.Request) {
	l := logutils.ZapFromRequest(r)

	s.isHealthyMx.Lock()
	if !s.isHealthy {
		s.isHealthyMx.Unlock()
		return
	}

	s.isHealthy = false
	l.Info("Server marked as unhealthy")

	// Let's not hold onto the lock in our sleep
	s.isHealthyMx.Unlock()

	// Give LB enough time to detect us unhealthy
	time.Sleep(s.cfg.DrainDuration)
}

func (s *Server) handleUndrain(w http.ResponseWriter, r *http.Request) {
	l := logutils.ZapFromRequest(r)

	s.isHealthyMx.Lock()
	defer s.isHealthyMx.Unlock()

	if s.isHealthy {
		return
	}

	s.isHealthy = true
	l.Info("Server marked as healthy")
}
