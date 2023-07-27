package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_Handlers_Healthcheck_Drain_Undrain(t *testing.T) {
	const (
		latency    = 200 * time.Millisecond
		listenAddr = ":8080"
	)

	//nolint: exhaustruct
	s := New(&Config{
		DrainDuration: latency,
		ListenAddr:    listenAddr,
		Log:           zap.Must(zap.NewDevelopment()),
	})

	{ // Check health
		req := httptest.NewRequest(http.MethodGet, "http://localhost"+listenAddr+"/health", nil)
		w := httptest.NewRecorder()
		s.handleHealthcheck(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Healthcheck must return 200 before draining")
	}

	{ // Drain
		req := httptest.NewRequest(http.MethodGet, "http://localhost"+listenAddr+"/drain", nil)
		w := httptest.NewRecorder()
		start := time.Now()
		s.handleDrain(w, req)
		duration := time.Since(start)
		resp := w.Result()
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Must return 200 for calls to `/drain`")
		assert.GreaterOrEqual(t, duration, latency, "Must wait long enough during draining")
	}

	{ // Check health
		req := httptest.NewRequest(http.MethodGet, "http://localhost"+listenAddr+"/health", nil)
		w := httptest.NewRecorder()
		s.handleHealthcheck(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode, "Healthcheck must return 500 after draining")
	}

	{ // Undrain
		req := httptest.NewRequest(http.MethodGet, "http://localhost"+listenAddr+"/undrain", nil)
		w := httptest.NewRecorder()
		s.handleUndrain(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Must return 200 for calls to `/undrain`")
		time.Sleep(latency)
	}

	{ // Check health
		req := httptest.NewRequest(http.MethodGet, "http://localhost"+listenAddr+"/health", nil)
		w := httptest.NewRecorder()
		s.handleHealthcheck(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Healthcheck must return 200 after undraining")
	}
}
