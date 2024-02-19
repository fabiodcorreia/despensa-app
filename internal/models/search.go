package models

type SearchResult struct {
	Id    string
	Label string
	Type  string
}

func NewSearchResult(id, label string) SearchResult {
	return SearchResult{
		Id:    id,
		Label: label,
		Type:  "location",
	}
}
