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
	// sess, _ := session.Get("session", ctx)
	// Need to check the error in case session is not there
	// sess.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   86400 * 7, // 1 week
	// 	HttpOnly: true,
	// 	Secure:   true,
	// }
	// sess.Values["foo"] = "bar"
	// sess.Save(ctx.Request(), ctx.Response())
	return render(ctx, pages.ViewHome("Despensa"))
}
