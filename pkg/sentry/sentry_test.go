package sentry

import (
	"testing"

	"github.com/harrisonlob/goyagi/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.Config{
		Environment: "test",
	}
	sentry, err := New(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, sentry)
}
