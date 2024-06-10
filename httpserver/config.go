package httpserver

import (
	"flag"
	"log/slog"
	"time"

	"github.com/flashbots/go-template/common"
	"github.com/flashbots/go-utils/envflag"
)

type flags struct {
	metricsAddr  *string
	drainSeconds *int
	listenAddr   *string

	logJSON    *bool
	logDebug   *bool
	logService *string
}

func defaults() flags {
	fg := flags{
		drainSeconds: envflag.MustInt("drain-seconds", 45, "seconds to wait in drain HTTP request"),
		listenAddr:   envflag.String("listen-addr", "127.0.0.1:8080", "address to listen on"),
		logJSON:      envflag.MustBool("log-json", false, "log in JSON format"),
		logDebug:     envflag.MustBool("log-debug", false, "log debug messages"),
		logService:   envflag.String("log-service", "your-project", "\"service\" tag to logs"),
		metricsAddr:  envflag.String("metrics-addr", "", "address to listen on for prometheus metrics"),
	}
	flag.Parse()
	return fg
}

// -----------------------------------------------------------------------------

type Config struct {
	DrainDuration            time.Duration
	GracefulShutdownDuration time.Duration
	ListenAddr               string
	Log                      *slog.Logger
	MetricsAddr              string
	ReadTimeout              time.Duration
	Version                  string
	WriteTimeout             time.Duration
}

func NewConfig(version string) *Config {
	flags := defaults()
	log := common.SetupLogger(&common.LoggingOpts{
		Debug:   *flags.logDebug,
		JSON:    *flags.logJSON,
		Service: *flags.logService,
		Version: version,
	})

	cfg := &Config{
		DrainDuration:            time.Duration(*flags.drainSeconds) * time.Second,
		GracefulShutdownDuration: 30 * time.Second,
		ListenAddr:               *flags.listenAddr,
		Log:                      log,
		MetricsAddr:              *flags.metricsAddr,
		ReadTimeout:              60 * time.Second,
		Version:                  version,
		WriteTimeout:             30 * time.Second,
	}

	if cfg.DrainDuration >= cfg.ReadTimeout {
		log.With("drainDuration", cfg.DrainDuration).With("readTimeout", cfg.ReadTimeout).Warn("Drain duration is not shorter that read timeout")
	}

	return cfg
}
