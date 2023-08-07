package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/flashbots/go-template/httpserver"
	"github.com/flashbots/go-utils/logutils"
)

var version = "dev" // is set during build process

func main() {
	cfg := httpserver.NewConfig(version)

	// Make sure to flush the logger before exiting the app
	defer logutils.FlushZap(cfg.Log)

	// Run server in background and wait for termination signal
	srv := httpserver.New(cfg)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	srv.RunInBackground()
	<-exit

	// Shutdown server once termination signal is received
	srv.Shutdown()
}
