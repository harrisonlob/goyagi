package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/binder"
	"github.com/harrisonlob/goyagi/pkg/errors"
	"github.com/harrisonlob/goyagi/pkg/health"
	"github.com/harrisonlob/goyagi/pkg/movies"
	"github.com/harrisonlob/goyagi/pkg/recovery"
	"github.com/harrisonlob/goyagi/pkg/signals"
	"github.com/labstack/echo"
	"github.com/lob/logger-go"
)

// New returns a new HTTP server with the registered routes.
func New(app application.App) *http.Server {
	log := logger.New()

	e := echo.New()

	b := binder.New()
	e.Binder = b

	e.Use(logger.Middleware())
	e.Use(recovery.Middleware())

	errors.RegisterErrorHandler(e, app)

	health.RegisterRoutes(e)
	movies.RegisterRoutes(e, app)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
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
