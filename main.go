package main

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"time"

	"github.com/flashbots/go-template/server"
)

var version = "dev" // is set during build process

func main() {
	cfg := server.NewConfig(version)

	// Make sure to flush the logger before exiting the app
	defer func() {
		if err := cfg.Log.Sync(); err != nil {
			// Workaround for `inappropriate ioctl for device` or `invalid argument` errors
			// See: https://github.com/uber-go/zap/issues/880#issuecomment-731261906
			var pathErr *fs.PathError
			if errors.As(err, &pathErr) {
				if pathErr.Path == "/dev/stderr" && pathErr.Op == "sync" {
					return
				}
			}
			fmt.Fprintf(
				os.Stderr,
				"{\"level\":\"error\",\"ts\":\"%s\",\"msg\":\"Failed to sync the logger\",\"error\":\"%s\"}\n",
				time.Now().Format(time.RFC3339),
				html.EscapeString(err.Error()),
			)
		}
	}()

	srv := server.New(cfg)

	srv.Run()
}
