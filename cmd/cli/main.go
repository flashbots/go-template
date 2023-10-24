package main

import (
	"flag"
	"strings"

	"github.com/flashbots/go-utils/envflag"
	"github.com/flashbots/go-utils/logutils"
	"go.uber.org/zap"
)

var (
	version = "dev" // is set during build process

	logProd    = envflag.MustBool("log-prod", false, "log in production mode (json)")
	logLevel   = envflag.String("log-level", "info", "log level (one of: \""+strings.Join(logutils.Levels, "\", \"")+"\")")
	logService = envflag.String("log-service", "your-project", "\"service\" tag to logs")

	log = logutils.MustGetZapLogger(
		logutils.LogDevMode(!*logProd),
		logutils.LogLevel(*logLevel),
	).With(zap.String("version", version))
)

func main() {
	flag.Parse()

	// Finish setting up logger, if needed
	if *logService != "" {
		log = log.With(zap.String("service", *logService))
	}
	defer logutils.FlushZap(log) // Makes sure that logger is flushed before the app exits

	log.Info("Starting the project")

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message (stacktrace added automatically when in log-dev mode)")
	log.Error("error message (stacktrace added automatically)")
	// log.Fatal("fatal message (stacktrace added automatically + causes the app to exit with non-zero status)")
}
