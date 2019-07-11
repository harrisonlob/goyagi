package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/getsentry/raven-go"
	"github.com/harrisonlob/goyagi/internal/test"
	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/sentry"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type mockSentryClient struct{}

func (m mockSentryClient) Capture(packet *raven.Packet, captureTags map[string]string) (string, chan error) {
	return "", nil
}

func TestHandler(t *testing.T) {
	app := application.App{
		Sentry: sentry.Sentry{
			Client: mockSentryClient{},
		},
	}
	h := handler{app}

	t.Run("surfaces generic errors as internal server errors", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON)
		err := errors.New("foo")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusInternalServerError, rr.Code, "expected generic errors to be 500s")
		assert.Contains(tt, rr.Body.String(), "Internal Server Error", "expected generic errors to have the correct message")
	})

	t.Run("surfaces HTTP errors transparently but obfuscates message", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON)
		err := echo.NewHTTPError(http.StatusForbidden, "foo")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusForbidden, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Forbidden", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 400 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON)
		err := echo.NewHTTPError(http.StatusBadRequest, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusBadRequest, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Bad Request", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 403 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON)
		err := echo.NewHTTPError(http.StatusForbidden, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusForbidden, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Forbidden", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 404 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON)
		err := echo.NewHTTPError(http.StatusNotFound, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusNotFound, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Not Found", "expected HTTP errors to have the correct message")
	})
}
