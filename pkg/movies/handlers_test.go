package movies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("lists movies on success", func(tt *testing.T) {
		c, rr := newContext(tt, nil)

		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, len(response) >= 23)
	})

	t.Run("Test limit and offset parameter", func(tt *testing.T) {
		offset := 2
		byteStr := fmt.Sprintf(`{"limit": 3, "offset": %d}`, offset)
		payload := []byte(byteStr)
		c, rr := newContext(tt, payload)

		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.Equal(tt, 3, len(response))
		// should be Captain Marvel because we first order in descending order by ID and then offset from that. Therefore,
		// in the 'add_movies' file Captain Marvel row is 3rd to last meaning it has the 3rd highest ID values and thus will
		// be the first value returned in this SQL query
		fullContext, fullRR := newContext(tt, nil)
		h.listHandler(fullContext)
		var fullResponse []model.Movie
		json.Unmarshal(fullRR.Body.Bytes(), &fullResponse)
		assert.True(tt, fullResponse[0].ID > response[0].ID)
	})

}

func TestRetrieveHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("retrieves movie on success", func(tt *testing.T) {
		c, rr := newContext(tt, nil)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.retrieveHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.Equal(tt, 1, response.ID)
		assert.Equal(tt, "Iron Man", response.Title)
	})

	t.Run("returns 404 if user isn't found", func(tt *testing.T) {
		c, _ := newContext(tt, nil)
		c.SetParamNames("id")
		c.SetParamValues("9999")

		err := h.retrieveHandler(c)
		assert.Contains(tt, err.Error(), "movie not found")
	})
}

func TestCreateHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("testing normal entry", func(tt *testing.T) {
		payload := []byte(`{"title": "Movie!"}`)
		c, rr := newContext(tt, payload)

		err := h.createHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(tt, "Movie!", response.Title)

		byteStr := fmt.Sprintf(`{"title": %q}`, response.Title)
		deletePayload := []byte(byteStr)
		deleteContext, deleteRR := newContext(tt, deletePayload)
		deleteErr := h.deleteHandler(deleteContext)
		assert.NoError(tt, deleteErr)
		assert.Equal(tt, http.StatusOK, deleteRR.Code)
	})

	t.Run("testing invalid no title entry", func(tt *testing.T) {
		payload := []byte(`{}`)
		c, _ := newContext(tt, payload)

		err := h.createHandler(c)
		assert.Error(tt, err)
	})
}

func newHandler(t *testing.T) handler {
	t.Helper()

	app, err := application.New()
	require.NoError(t, err)
	return handler{app}
}

func newContext(t *testing.T, payload []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)
	return c, rr
}
