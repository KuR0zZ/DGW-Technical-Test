package repository

import (
	"database/sql"
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
	query := "UPDATE Books SET name = $1, genre = $2, author = $3, published_date = $4, stock = $5, price = $6 WHERE id = $7"

	_, err := repository.DB.Exec(query, book.Name, book.Genre, book.Author, book.PublishedDate, book.Stock, book.Price, book.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *BookRepositoryImpl) Delete(bookId int) error {
	query := "DELETE FROM Books WHERE id = $1"

	result, err := repository.DB.Exec(query, bookId)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (repository *BookRepositoryImpl) FindAll() ([]entity.Book, error) {
	query := "SELECT * FROM Books"

	var books []entity.Book
	if err := repository.DB.Select(&books, query); err != nil {
		return nil, err
	}

	return books, nil
}

func (repository *BookRepositoryImpl) FindById(bookId int) (*entity.Book, error) {
	query := "SELECT * FROM Books WHERE id = $1"

	book := new(entity.Book)
	if err := repository.DB.Get(book, query, bookId); err != nil {
		return nil, err
	}

	return book, nil
}
