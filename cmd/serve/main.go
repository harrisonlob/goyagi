package main

import (
	"net/http"

	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/server"
	"github.com/lob/logger-go"
)

func main() {
	log := logger.New()

	app, err := application.New()
	if err != nil {
		log.Err(err).Fatal("failed to initialize application")
	}

	srv := server.New(app)

	log.Info("server started", logger.Data{"port": app.Config.Port})

	srv_err := srv.ListenAndServe()
	if srv_err != nil && srv_err != http.ErrServerClosed {
		log.Err(srv_err).Fatal("server stopped")
	}

	log.Info("server stopped")
}
