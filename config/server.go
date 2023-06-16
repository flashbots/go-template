// Package config provides basic primitives for server configuration
package config

import (
	"flag"
	"time"

	"github.com/flashbots/go-template/util"
	"go.uber.org/zap"
)

type flags struct {
	debug        *bool
	drainSeconds *int
	listenAddr   *string
	logProd      *bool
	logService   *string
}

func defaults() flags {
	fg := flags{
		debug:        util.FlagB("debug", false, "print debug output"),
		drainSeconds: util.FlagI("drain-seconds", 45, "seconds to wait in drain HTTP request"),
		listenAddr:   util.FlagS("listen-addr", "127.0.0.1:8080", "address to listen on"),
		logProd:      util.FlagB("log-prod", true, "log in production mode (json)"),
		logService:   util.FlagS("log-service", "your-project", "'service' tag to logs"),
	}
	flag.Parse()
	return fg
}

// -----------------------------------------------------------------------------

type Server struct {
	DrainDuration            time.Duration
	GracefulShutdownDuration time.Duration
	ListenAddr               string
	Log                      *zap.Logger
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	Version                  string
}

func NewServerConfig(version string) *Server {
	flags := defaults()
	log := getLogger(flags)

	srv := &Server{
		DrainDuration:            time.Duration(*flags.drainSeconds) * time.Second,
		GracefulShutdownDuration: 30 * time.Second,
		ListenAddr:               *flags.listenAddr,
		Log:                      log,
		ReadTimeout:              60 * time.Second,
		Version:                  version,
		WriteTimeout:             30 * time.Second,
	}

	if srv.DrainDuration >= srv.ReadTimeout {
		log.Warn("Drain duration is not shorter that read timeout",
			zap.Duration("drainDuration", srv.DrainDuration),
			zap.Duration("readTimeout", srv.ReadTimeout),
		)
	}

	return srv
}
