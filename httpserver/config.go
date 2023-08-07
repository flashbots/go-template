package httpserver

import (
	"flag"
	"strings"
	"time"

	"github.com/flashbots/go-utils/envflag"
	"github.com/flashbots/go-utils/logutils"
	"go.uber.org/zap"
)

type flags struct {
	drainSeconds *int
	listenAddr   *string
	logDev       *bool
	logLevel     *string
	logService   *string
}

func defaults() flags {
	fg := flags{
		drainSeconds: envflag.MustInt("drain-seconds", 45, "seconds to wait in drain HTTP request"),
		listenAddr:   envflag.String("listen-addr", "127.0.0.1:8080", "address to listen on"),
		logDev:       envflag.MustBool("log-dev", false, "log in development mode (plain text)"),
		logLevel:     envflag.String("log-level", "info", "log level (one of: \""+strings.Join(logutils.Levels, "\", \"")+"\")"),
		logService:   envflag.String("log-service", "your-project", "\"service\" tag to logs"),
	}
	flag.Parse()
	return fg
}

// -----------------------------------------------------------------------------

type Config struct {
	DrainDuration            time.Duration
	GracefulShutdownDuration time.Duration
	ListenAddr               string
	Log                      *zap.Logger
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	Version                  string
}

func NewConfig(version string) *Config {
	flags := defaults()
	log := logutils.MustGetZapLogger(
		logutils.LogDevMode(*flags.logDev),
		logutils.LogLevel(*flags.logLevel),
	).With(
		zap.String("service", *flags.logService),
		zap.String("version", version),
	)

	cfg := &Config{
		DrainDuration:            time.Duration(*flags.drainSeconds) * time.Second,
		GracefulShutdownDuration: 30 * time.Second,
		ListenAddr:               *flags.listenAddr,
		Log:                      log,
		ReadTimeout:              60 * time.Second,
		Version:                  version,
		WriteTimeout:             30 * time.Second,
	}

	if cfg.DrainDuration >= cfg.ReadTimeout {
		log.Warn("Drain duration is not shorter that read timeout",
			zap.Duration("drainDuration", cfg.DrainDuration),
			zap.Duration("readTimeout", cfg.ReadTimeout),
		)
	}

	return cfg
}
