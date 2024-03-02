package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	sv services.ItemService
}

func NewItemHandler(sv services.ItemService) ItemHandler {
	return ItemHandler{
		sv,
	}
}

func (h ItemHandler) View(ctx echo.Context) error {
	itemID := getParam(ctx, KeyID)

	loc, err := h.sv.GetItem(itemID)
	if err != nil {
		if errors.Is(err, services.ErrItemNotFound) {
			return renderNotFound(ctx)
		}
		return err
	}
	return ctx.JSON(http.StatusOK, loc)
}

func (h ItemHandler) List(ctx echo.Context) error {
	items, err := h.sv.GetItems()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, items)
}

func (h ItemHandler) Add(ctx echo.Context) error {
	quantity, _ := strconv.ParseInt(getFormValue(ctx, "quantity"), 10, 8)
	item := models.NewItem(getFormValue(ctx, KeyName))
	item.LocationID = getFormValue(ctx, "locationId")
	item.Quantity = int8(quantity)

	err := h.sv.AddItem(item)
	if err != nil {
		return err
	}

	return renderEmpty(ctx)
}
