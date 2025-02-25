package repository

import (
	"dgw-technical-test/entity"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	Create(book *entity.Book) error
	Update(book *entity.Book) error
	Delete(bookId int) error
	FindAll() ([]entity.Book, error)
	FindById(bookId int) (*entity.Book, error)
}

type BookRepositoryImpl struct {
	DB *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepositoryImpl {
	return &BookRepositoryImpl{DB: db}
}

func (repository *BookRepositoryImpl) Create(book *entity.Book) error {
	query := "INSERT INTO Books (name, genre, author, published_date, stock, price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	if err := repository.DB.QueryRow(query, book.Name, book.Genre, book.Author, book.PublishedDate, book.Stock, book.Price).Scan(&book.ID); err != nil {
		return err
	}

	return nil
}

func (repository *BookRepositoryImpl) Update(book *entity.Book) error {
	return nil
}

func (repository *BookRepositoryImpl) Delete(bookId int) error {
	return nil
}

func (repository *BookRepositoryImpl) FindAll() ([]entity.Book, error) {
	return nil, nil
}

func (repository *BookRepositoryImpl) FindById(bookId int) (*entity.Book, error) {
	return nil, nil
}
