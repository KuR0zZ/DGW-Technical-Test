package dto

type BookCreateRequest struct {
	Name          string  `json:"name" validate:"required"`
	Genre         string  `json:"genre" validate:"required"`
	Author        string  `json:"author" validate:"required"`
	PublishedDate string  `json:"published_date" validate:"required"`
	Stock         int     `json:"stock" validate:"required"`
	Price         float64 `json:"price" validate:"required"`
}

type BookCreateResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Genre         string  `json:"genre"`
	Author        string  `json:"author"`
	PublishedDate string  `json:"published_date"`
	Stock         int     `json:"stock"`
	Price         float64 `json:"price"`
}

type BookUpdateRequest struct {
	Name          string  `json:"name" validate:"required"`
	Genre         string  `json:"genre" validate:"required"`
	Author        string  `json:"author" validate:"required"`
	PublishedDate string  `json:"published_date" validate:"required"`
	Stock         int     `json:"stock" validate:"required"`
	Price         float64 `json:"price" validate:"required"`
}
