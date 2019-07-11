package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/harrisonlob/goyagi/pkg/config"
	"github.com/pkg/errors"
)

type ravenClient interface {
	Capture(packet *raven.Packet, captureTags map[string]string) (string, chan error)
}

type Sentry struct {
	Client ravenClient
}

func New(cfg config.Config) (Sentry, error) {
	defaultTags := map[string]string{
		"environment": cfg.Environment,
	}
	client, err := raven.NewWithTags(cfg.SentryDSN, defaultTags)
	if err != nil {
		return Sentry{}, errors.Wrap(err, "sentry")
	}

	return Sentry{client}, nil
}
