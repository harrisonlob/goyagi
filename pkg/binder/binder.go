package binder

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/creasty/defaults"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/mold.v2"
	"gopkg.in/go-playground/mold.v2/modifiers"
	"gopkg.in/go-playground/validator.v9"
)

type Binder struct {
	binder   *echo.DefaultBinder
	conform  *mold.Transformer
	validate *validator.Validate
}

func New() *Binder {
	binder := &echo.DefaultBinder{}
	conform := modifiers.New()
	validate := validator.New()

	return &Binder{binder, conform, validate}
}

func (b *Binder) Bind(i interface{}, c echo.Context) error {
	if err := b.binder.Bind(i, c); err != nil {
		return err
	}

	if err := b.conform.Struct(context.Background(), i); err != nil {
		return err
	}

	if err := defaults.Set(i); err != nil {
		return err
	}

	if err := b.validate.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		msg := format(errs[0])
		return echo.NewHTTPError(http.StatusUnprocessableEntity, msg)
	}

	return nil
}

func format(err validator.FieldError) string {
	if err.Kind() == reflect.Int {
		switch err.Tag() {
		case "max":
			return fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
		case "min":
			return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
		}
	} else if err.Kind() == reflect.String {
		switch err.Tag() {
		case "max":
			return fmt.Sprintf("%s length must be less than or equal to %s characters long", err.Field(), err.Param())
		case "min":
			return fmt.Sprintf("%s length must be at least %s characters long", err.Field(), err.Param())
		}
	}

	if err.Tag() == "required" {
		return fmt.Sprintf("%s is required", err.Field())
	}

	return err.Param()
}
