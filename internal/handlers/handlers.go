package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
)

const KeyID = `id`
const KeyName = `name`

var policy = bluemonday.StrictPolicy()

func getParam(ctx echo.Context, name string) string {
	return policy.Sanitize(ctx.Param(name))
}

func getFormValue(ctx echo.Context, name string) string {
	return policy.Sanitize(ctx.FormValue(name))
}

func render(ctx echo.Context, comp templ.Component) error {
	return renderWithStatusCode(ctx, comp, http.StatusOK)
}

func renderNotFound(ctx echo.Context) error {
	return renderWithStatusCode(ctx, templ.NopComponent, http.StatusNotFound)
}

func renderEmpty(ctx echo.Context) error {
	return renderWithStatusCode(ctx, templ.NopComponent, http.StatusOK)
}

func renderWithStatusCode(ctx echo.Context, comp templ.Component, statusCode int) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return comp.Render(ctx.Request().Context(), ctx.Response().Writer)
}
