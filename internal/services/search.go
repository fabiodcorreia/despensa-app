package services

import (
	"fmt"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

type SearchService struct {
	store *storage.Store
}

func NewSearch(s *storage.Store) *SearchService {
	return &SearchService{
		store: s,
	}
}

func (s *SearchService) FindAll(filter string) ([]models.SearchResult, error) {
	locs, err := s.store.FilterLocations(filter)
	if err != nil {
		return nil, fmt.Errorf("search service find all fail: %w", err)
	}

	results := make([]models.SearchResult, 0, len(locs))
	for _, loc := range locs {
		results = append(results, models.NewSearchResult(loc.Id, loc.Name))
	}

	return results, nil
}
