package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/harrisonlob/goyagi/pkg/health"
	"github.com/harrisonlob/goyagi/pkg/signals"
	"github.com/labstack/echo"
	"github.com/lob/logger-go"
)

// New returns a new HTTP server with the registered routes.
func New() *http.Server {
	log := logger.New()

	e := echo.New()

	health.RegisterRoutes(e)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 3000),
		Handler: e,
	}

	graceful := signals.Setup()

	go func() {
		<-graceful
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Err(err).Error("server shutdown")
		}
	}()

	return srv
}
