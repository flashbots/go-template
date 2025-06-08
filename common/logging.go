// Package common contains common utilities and functions used by the service.
package common

import (
	"log/slog"

	"github.com/go-chi/httplog/v2"
)

type LoggingOpts struct {
	Service        string
	JSON           bool
	Debug          bool
	Concise        bool
	RequestHeaders bool
	Version        string
	UID            string
	Tags           map[string]string
}

func SetupLogger(opts *LoggingOpts) (log *httplog.Logger) {
	logLevel := slog.LevelInfo
	if opts.Debug {
		logLevel = slog.LevelDebug
	}

	// If version is provided, add it to the tags.
	if opts.Version != "" || opts.UID != "" {
		if opts.Tags == nil {
			opts.Tags = make(map[string]string)
		}
		if opts.Version != "" {
			opts.Tags["version"] = opts.Version
		}
		if opts.UID != "" {
			opts.Tags["uid"] = opts.UID
		}
	}

	logger := httplog.NewLogger(opts.Service, httplog.Options{
		JSON:           opts.JSON,
		LogLevel:       logLevel,
		Concise:        opts.Concise,
		RequestHeaders: opts.RequestHeaders,
		Tags:           opts.Tags,
	})

	return logger
}
