package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		// he.Message
	}

	if code == http.StatusNotFound {
		if err := renderNotFound(ctx); err != nil {
			ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
		}
		return
	}

	ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
}
