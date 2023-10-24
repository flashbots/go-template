package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/flashbots/go-utils/envflag"
)

var (
	version = "dev" // is set during build process

	logProd    = envflag.MustBool("log-prod", false, "log in production mode (json)")
	logDebug   = envflag.MustBool("log-debug", false, "log debug messages")
	logService = envflag.String("log-service", "", "'service' tag to logs")
)

type LoggingOpts struct {
	Debug   bool
	JSON    bool
	Service string
	Version string
}

func setupLogger(opts *LoggingOpts) (log *slog.Logger) {
	logLevel := slog.LevelInfo
	if opts.Debug {
		logLevel = slog.LevelDebug
	}

	if opts.JSON {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	} else {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	if opts.Service != "" {
		log = log.With("service", opts.Service)
	}

	if opts.Version != "" {
		log = log.With("version", opts.Version)
	}

	return log
}

func main() {
	flag.Parse()
	log := setupLogger(&LoggingOpts{*logDebug, *logProd, *logService, version})
	log.Info("Starting the project")

	log.Debug("debug message")
	log.Info("info message")
	log.With("key", "value").Warn("warn message")
	log.Error("error message (stacktrace added automatically)")
	// log.Fatal("fatal message (stacktrace added automatically + causes the app to exit with non-zero status)")
}
