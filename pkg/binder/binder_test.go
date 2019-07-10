package binder

import (
	"testing"

	"github.com/harrisonlob/goyagi/internal/test"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type params struct {
	Title string `json:"title" mod:"trim" validate:"required,min=5,max=8"`
	Age   int    `json:"age" default:"21" validate:"min=18,max=99"`
}

func TestNew(t *testing.T) {
	b := New()
	assert.NotNil(t, b)
	assert.NotNil(t, b.binder)
	assert.NotNil(t, b.conform)
	assert.NotNil(t, b.validate)
}

func TestBind(t *testing.T) {
	b := New()
	assert.NotNil(t, b)

	t.Run("rejects invalid content types", func(tt *testing.T) {
		payload := []byte("{}")
		c, _ := test.NewContext(t, payload, "invalid")
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "Unsupported Media Type")
	})

	t.Run("extracts request from request", func(tt *testing.T) {
		payload := []byte(`{"title": "banana"}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.NoError(t, err)
		assert.Equal(t, "banana", p.Title)
	})

	t.Run("enforces required values", func(tt *testing.T) {
		payload := []byte(`{}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "is required")
	})

	t.Run("sets default values", func(tt *testing.T) {
		payload := []byte(`{"title": "banana"}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.NoError(t, err)
		assert.Equal(t, p.Age, 21)
	})

	t.Run("trims whitespace", func(tt *testing.T) {
		payload := []byte(`{"title": "       banana       "}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.NoError(t, err)
		assert.Equal(t, "banana", p.Title)
	})

	t.Run("enforces min values on strings", func(tt *testing.T) {
		payload := []byte(`{"title": "1"}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "length must be at least 5 characters long")
	})

	t.Run("enforces max values on strings", func(tt *testing.T) {
		payload := []byte(`{"title": "123456789"}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "length must be less than or equal to 8 characters long")
	})

	t.Run("enforces min values on integers", func(tt *testing.T) {
		payload := []byte(`{"title": "banana", "age": 1}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "must be at least 18")
	})

	t.Run("enforces max values on integers", func(tt *testing.T) {
		payload := []byte(`{"title": "banana", "age": 100}`)
		c, _ := test.NewContext(t, payload, echo.MIMEApplicationJSON)
		p := params{}
		err := b.Bind(&p, c)
		assert.Contains(t, err.Error(), "must be less than or equal to 99")
	})
}
