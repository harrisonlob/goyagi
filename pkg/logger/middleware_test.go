package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harrisonlob/goyagi/internal/test"
	"github.com/labstack/echo"
	"github.com/lob/logger-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(Middleware())

	e.GET("/", func(c echo.Context) error {
		log, ok := c.Get("logger").(logger.Logger)
		assert.True(t, ok, "expected logger to be of type Logger")
		assert.NotNil(t, log)
		return nil
	})

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	e.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFromContext(t *testing.T) {
	log := logger.New()
	c, _ := test.NewContext(t, nil, echo.MIMEApplicationJSON)
	c.Set(key, log)

	l := FromContext(c)

	assert.Equal(t, log, l)
}
