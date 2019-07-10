package movies

import "time"

type createParams struct {
	Title       string    `json:"title" mod:"trim" validate:"required,max=35"`
	ReleaseDate time.Time `json:"release_date" validate:"omitempty"`
}

type deleteParams struct {
	Title string `json:"title" validate:"required"`
}

type listParams struct {
	Limit  int `query:"limit" default:"10" validate:"min=0,max=100"`
	Offset int `query:"offset" validate:"min=0"`
}
