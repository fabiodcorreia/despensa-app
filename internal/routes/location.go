package routes

import (
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
)

const routeLocationRoot = "/location"

type LocationRoute struct {
	handler handlers.LocationHandler
}

var _ Routes = (*LocationRoute)(nil)

func NewLocationRoute(handler handlers.LocationHandler) LocationRoute {
	return LocationRoute{
		handler,
	}
}

func (r LocationRoute) Routes() []Route {
	return []Route{
		NewRoute(http.MethodGet, routeLocationRoot, r.handler.List),
		NewRoute(http.MethodPost, routeLocationRoot, r.handler.Add),
		NewRoute(http.MethodGet, routeLocationRoot+"/:"+handlers.KeyID, r.handler.View),
	}
}
