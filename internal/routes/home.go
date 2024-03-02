package routes

import (
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/handlers"
)

type HomeRoute struct {
	handler handlers.HomeHandler
}

var _ Routes = (*HomeRoute)(nil)

func (r HomeRoute) Routes() []Route {
	return []Route{
		NewRoute(http.MethodGet, "/", r.handler.View),
	}
}
