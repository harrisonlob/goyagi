package test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func NewContext(t *testing.T, payload []byte, mime string) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, mime)
	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)
	return c, rr
}
