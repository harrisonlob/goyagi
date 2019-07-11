package errors

import (
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/logger"
	"github.com/labstack/echo"
	loggergo "github.com/lob/logger-go"
)

type handler struct {
	app application.App
}

func RegisterErrorHandler(e *echo.Echo, app application.App) {
	h := handler{app}

	e.HTTPErrorHandler = h.handleError
}

func (h *handler) handleError(err error, c echo.Context) {
	log := logger.FromContext(c)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = http.StatusText(code)
	}

	if code == http.StatusInternalServerError {
		stacktrace := raven.NewException(err, raven.GetOrNewStacktrace(err, 0, 2, nil))
		httpContext := raven.NewHttp(c.Request())
		packet := raven.NewPacket(msg, stacktrace, httpContext)

		h.app.Sentry.Client.Capture(packet, map[string]string{})
	}

	log.Root(loggergo.Data{"status_code": code}).Err(err).Error("request error")

	err = c.JSON(code, map[string]interface{}{"error": map[string]interface{}{"message": msg, "status_code": code}})
	if err != nil {
		log.Err(err).Error("error handler json error")
	}
}
