package entity

import "time"

type Book struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	Genre         string    `db:"genre"`
	Author        string    `db:"author"`
	PublishedDate string    `db:"published_date"`
	Stock         int       `db:"stock"`
	Price         float64   `db:"price"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
