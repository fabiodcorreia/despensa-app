package handlers

import (
	"errors"
	"net/http"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/fabiodcorreia/despensa-app/internal/views/pages"
	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	ls services.LocationService
}

func NewLocationHandler(ls services.LocationService) LocationHandler {
	return LocationHandler{ls}
}

func (h LocationHandler) View(ctx echo.Context) error {
	locationID := getParam(ctx, KeyID)

	loc, err := h.ls.GetLocation(locationID)
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			return renderNotFound(ctx)
		}
		return err
	}
	//TODO: Render Component
	return ctx.JSON(http.StatusOK, loc)
}

func (h LocationHandler) List(ctx echo.Context) error {
	locations, err := h.ls.GetLocations()
	if err != nil {
		return err
	}
	return render(ctx, pages.ViewLocations("Despensa", locations))
}

func (h LocationHandler) Add(ctx echo.Context) error {
	name := getFormValue(ctx, KeyName)
	loc := models.NewLocation(name)

	err := h.ls.AddLocation(loc)
	if err != nil {
		if !errors.Is(err, services.ErrLocationExists) {
			return err
		}
	}

	return renderEmpty(ctx)
}
