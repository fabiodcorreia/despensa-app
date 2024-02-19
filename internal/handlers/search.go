package handlers

import (
	"time"

	"github.com/fabiodcorreia/despensa-app/internal/services"
	"github.com/fabiodcorreia/despensa-app/internal/views"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	s *services.SearchService
}

func NewSearch(s *services.SearchService) *SearchHandler {
	return &SearchHandler{
		s,
	}
}

func (h *SearchHandler) View(ctx echo.Context) error {
	return render(ctx, views.SearchView())
}

func (h *SearchHandler) Search(ctx echo.Context) error {
	searchText := ctx.FormValue("search")
	time.Sleep(2 * time.Second)
	results, err := h.s.FindAll(searchText)
	if err != nil {
		return err
	}

	return render(ctx, views.Result(results))
}
