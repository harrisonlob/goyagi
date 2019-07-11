package recovery

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := errors.Errorf("%v", r)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
