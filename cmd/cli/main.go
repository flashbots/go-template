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

	logDev     = envflag.MustBool("log-dev", false, "log in development mode (plain text)")
	logLevel   = envflag.String("log-level", "info", "log level (one of: \""+strings.Join(logutils.Levels, "\", \"")+"\")")
	logService = envflag.String("log-service", "your-project", "\"service\" tag to logs")
)

func main() {
	flag.Parse()

	// Setup logging
	l := logutils.MustGetZapLogger(
		logutils.LogDevMode(*logDev),
		logutils.LogLevel(*logLevel),
	).With(
		zap.String("service", *logService),
		zap.String("version", version),
	)
	defer logutils.FlushZap(l) // Makes sure that logger is flushed before the app exits

	l.Info("Starting the project")

	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message (stacktrace added automatically when in log-dev mode)")
	l.Error("error message (stacktrace added automatically)")
	// l.Fatal("fatal message (stacktrace added automatically + causes the app to exit with non-zero status)")
}
