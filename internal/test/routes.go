package test

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
)

func FilterRoutes(routes []*echo.Route) []string {
	rs := []string{}

	for _, r := range routes {
		if strings.Contains(r.Name, "*handler") {
			rs = append(rs, fmt.Sprintf("%s %s", r.Method, r.Path))
		}
	}

	return rs
}
