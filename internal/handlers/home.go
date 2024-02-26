package handlers

import (
	"github.com/fabiodcorreia/despensa-app/internal/views/pages"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func NewHome() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) View(ctx echo.Context) error {
	return render(ctx, pages.ViewHome("HOME VIEW"))
}
