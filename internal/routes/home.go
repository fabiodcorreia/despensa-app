package routes

import (
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
)

type Home struct {
	handler handlers.HomeHandler
}

var _ Routes = (*Home)(nil)

// Make this interface?
func (r *Home) Routes() []Route {
	return []Route{
		NewRoute(http.MethodGet, "/", r.handler.View),
	}
}
