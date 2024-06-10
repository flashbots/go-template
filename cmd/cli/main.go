package main

import (
	"errors"
	"flag"

	"github.com/flashbots/go-template/common"
	"github.com/flashbots/go-utils/envflag"
)

var (
	version = "dev" // is set during build process

	logProd    = envflag.MustBool("log-prod", false, "log in production mode (json)")
	logDebug   = envflag.MustBool("log-debug", false, "log debug messages")
	logService = envflag.String("log-service", "", "'service' tag to logs")
)

func main() {
	flag.Parse()
	log := common.SetupLogger(&common.LoggingOpts{
		Debug:   *logDebug,
		JSON:    *logProd,
		Service: *logService,
		Version: version,
	})
	log.Info("Starting the project")

	log.Debug("debug message")
	log.Info("info message")
	log.With("key", "value").Warn("warn message")

	log.Error("error message", "err", errors.ErrUnsupported)
	// log.Fatal("fatal message (causes the app to exit with non-zero status)")
}
