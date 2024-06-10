// Package common contains common utilities and functions used by the service.
package common

import (
	"log/slog"
	"os"
)

type LoggingOpts struct {
	Debug   bool
	JSON    bool
	Service string
	Version string
}

func SetupLogger(opts *LoggingOpts) (log *slog.Logger) {
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
