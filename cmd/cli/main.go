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

	log *slog.Logger
)

func main() {
	flag.Parse()

	logLevel := slog.LevelInfo
	if *logDebug {
		logLevel = slog.LevelDebug
	}
	if *logProd {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	} else {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}
	log = log.With("version", version)
	if *logService != "" {
		log = log.With("service", *logService)
	}
	if *logProd {
		log = log.With("env", "prod")
	}

	log.Info("Starting the project")

	log.Debug("debug message")
	log.Info("info message")
	log.With("key", "value").Warn("warn message")
	log.Error("error message (stacktrace added automatically)")
	// log.Fatal("fatal message (stacktrace added automatically + causes the app to exit with non-zero status)")
}
