package application

import (
	"github.com/go-pg/pg"
	"github.com/harrisonlob/goyagi/pkg/config"
	"github.com/harrisonlob/goyagi/pkg/database"
	"github.com/harrisonlob/goyagi/pkg/sentry"
	"github.com/pkg/errors"
)

type App struct {
	Config config.Config
	DB     *pg.DB
	Sentry sentry.Sentry
}

func New() (App, error) {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	sentry, err := sentry.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	return App{cfg, db, sentry}, nil
}
