package logger

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/lob/logger-go"
	"github.com/pkg/errors"
)

const key = "logger"

func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	l := logger.New()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t1 := time.Now()

			id, err := uuid.NewV4()
			if err != nil {
				return errors.WithStack(err)
			}

			log := l.ID(id.String())
			c.Set(key, log)

			if err := next(c); err != nil {
				c.Error(err)
			}

			t2 := time.Now()

			var ipAddress string
			if xff := c.Request().Header.Get("x-forwarded-for"); xff != "" {
				split := strings.Split(xff, ",")
				ipAddress = strings.TrimSpace(split[len(split)-1])
			} else {
				ipAddress = c.Request().RemoteAddr
			}

			log.Root(logger.Data{
				"status_code":   c.Response().Status,
				"method":        c.Request().Method,
				"path":          c.Request().URL.Path,
				"route":         c.Path(),
				"response_time": t2.Sub(t1).Seconds() * 1000,
				"referer":       c.Request().Referer(),
				"user_agent":    c.Request().UserAgent(),
				"ip_address":    ipAddress,
				"trace_id":      c.Request().Header.Get("x-amzn-trace-id"),
			}).Info("request handled")

			return nil
		}
	}
}

func FromContext(c echo.Context) logger.Logger {
	if log, ok := c.Get(key).(logger.Logger); ok {
		return log
	}

	return logger.New()
}
