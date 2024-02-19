package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(ctx echo.Context, comp templ.Component) error {
	return renderWithStatusCode(ctx, comp, http.StatusOK)
}

func renderEmpty(ctx echo.Context) error {
	// I would like to use 204 No Content but HTMX only works with 200 as far I can find
	return renderWithStatusCode(ctx, templ.NopComponent, http.StatusOK)
}

func renderWithStatusCode(ctx echo.Context, comp templ.Component, statusCode int) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return comp.Render(ctx.Request().Context(), ctx.Response().Writer)
}
