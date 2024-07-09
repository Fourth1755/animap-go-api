package dtos

type AnimeDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Episodes int    `json:"episodes"`
	Seasonal string `json:"seasonal"`
	Year     string `json:"year"`
}

type AnimeQueryDTO struct {
	Seasonal string
	Year     string
}
