package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flashbots/go-template/common"
	"github.com/flashbots/go-template/httpserver"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:  "listen-addr",
		Value: "127.0.0.1:8080",
		Usage: "address to listen on for API",
	},
	&cli.StringFlag{
		Name:  "metrics-addr",
		Value: "127.0.0.1:8090",
		Usage: "address to listen on for Prometheus metrics",
	},
	&cli.BoolFlag{
		Name:  "log-json",
		Value: false,
		Usage: "log in JSON format",
	},
	&cli.BoolFlag{
		Name:  "log-debug",
		Value: false,
		Usage: "log debug messages",
	},
	&cli.StringFlag{
		Name:  "log-service",
		Value: "your-project",
		Usage: "add 'service' tag to logs",
	},
	&cli.Int64Flag{
		Name:  "drain-seconds",
		Value: 45,
		Usage: "seconds to wait in drain HTTP request",
	},
}

func main() {
	app := &cli.App{
		Name:  "httpserver",
		Usage: "Serve API, and metrics",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			listenAddr := cCtx.String("listen-addr")
			metricsAddr := cCtx.String("metrics-addr")
			logJSON := cCtx.Bool("log-json")
			logDebug := cCtx.Bool("log-debug")
			logService := cCtx.String("log-service")
			drainDuration := time.Duration(cCtx.Int64("drain-seconds")) * time.Second

			log := common.SetupLogger(&common.LoggingOpts{
				Debug:   logDebug,
				JSON:    logJSON,
				Service: logService,
				Version: common.Version,
			})

			cfg := &httpserver.HTTPServerConfig{
				ListenAddr:  listenAddr,
				MetricsAddr: metricsAddr,
				Log:         log,

				DrainDuration:            drainDuration,
				GracefulShutdownDuration: 30 * time.Second,
				ReadTimeout:              60 * time.Second,
				WriteTimeout:             30 * time.Second,
			}

			srv, err := httpserver.New(cfg)
			if err != nil {
				cfg.Log.Error("failed to create server", "err", err)
				return err
			}

			exit := make(chan os.Signal, 1)
			signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
			srv.RunInBackground()
			<-exit

			// Shutdown server once termination signal is received
			srv.Shutdown()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
