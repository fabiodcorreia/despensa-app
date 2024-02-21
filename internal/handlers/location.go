package handlers

import (
	"net/http"
	"strings"

	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/fabiodcorreia/despensa-app/internal/views"
	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	s *services.LocationService
}

func NewLocation(s *services.LocationService) *LocationHandler {
	return &LocationHandler{
		s,
	}
}

func (h *LocationHandler) View(ctx echo.Context) error {
	locationId := ctx.Param("id")
	locationId = strings.Trim(locationId, " ")
	if locationId == "" {
		return ctx.Redirect(http.StatusPermanentRedirect, "/search")
	}

	// Get the location from the DB
	loc, err := h.s.GetLocation(locationId)
	if err != nil {
		// If not found return 404
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Get the products on that items
	items, err := h.s.GetLocationItems(loc.Id)
	if err != nil {
		return renderEmpty(ctx)
	}

	// Render the list and return
	return render(ctx, views.LocationView(loc, items))
}
