package dtos

import "github.com/google/uuid"

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

type SearchResultItem struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
	Type  string    `json:"type"`
}

type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
}
