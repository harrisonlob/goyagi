package movies

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/harrisonlob/goyagi/pkg/application"
	"github.com/harrisonlob/goyagi/pkg/model"
	"github.com/labstack/echo"
)

type handler struct {
	app application.App
}

func (h *handler) listHandler(c echo.Context) error {
	params := listParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	var movies []*model.Movie

	err := h.app.DB.
		Model(&movies).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("id DESC").
		Select()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

func (h *handler) retrieveHandler(c echo.Context) error {
	id := c.Param("id")

	var movie model.Movie

	err := h.app.DB.Model(&movie).Where("id = ?", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "movie not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h *handler) createHandler(c echo.Context) error {
	params := createParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	movie := model.Movie{
		Title:       params.Title,
		ReleaseDate: params.ReleaseDate,
	}

	insertTimer := h.app.Metrics.NewTimer("goyagi.movies.create.db")

	_, err := h.app.DB.Model(&movie).Insert()
	if err != nil {
		insertTimer.End("result:error")
		return err
	}

	insertTimer.End("result:success")
	return c.JSON(http.StatusOK, movie)
}

func (h *handler) deleteHandler(c echo.Context) error {
	deletionParams := deleteParams{}
	if err := c.Bind(&deletionParams); err != nil {
		return err
	}

	movie := model.Movie{
		Title: deletionParams.Title,
	}
	res, err := h.app.DB.Model(&movie).Where("title = ?", movie.Title).Delete()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
