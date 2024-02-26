package routes

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	Method string
	Path   string
	Handle echo.HandlerFunc
}

// Create a new route with the given method, path and handler.
//
// method: http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,...
func NewRoute(method, path string, handle echo.HandlerFunc) Route {
	return Route{
		Method: method,
		Path:   path,
		Handle: handle,
	}
}

type Routes interface {
	Routes() []Route
}
