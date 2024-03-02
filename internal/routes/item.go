package routes

import (
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
)

const routeItemRoot = "/item"

type ItemRoute struct {
	handler handlers.ItemHandler
}

func NewItemRoute(handler handlers.ItemHandler) ItemRoute {
	return ItemRoute{
		handler,
	}
}

var _ Routes = (*ItemRoute)(nil)

func (r ItemRoute) Routes() []Route {
	return []Route{
		NewRoute(http.MethodGet, routeItemRoot, r.handler.List),
		NewRoute(http.MethodPost, routeItemRoot, r.handler.Add),
		NewRoute(http.MethodGet, routeItemRoot+"/:"+handlers.KeyID, r.handler.View),
		// NewRoute(http.MethodPut, "/item/:id"),
	}
}
