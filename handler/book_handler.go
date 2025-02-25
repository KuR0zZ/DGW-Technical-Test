package handler

import (
	"database/sql"
	"dgw-technical-test/dto"
	"dgw-technical-test/entity"
	"dgw-technical-test/repository"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BookHandler struct {
	BookRepository repository.BookRepository
	Validate       *validator.Validate
}

func NewBookHandler(bookRepository repository.BookRepository, validate *validator.Validate) *BookHandler {
	return &BookHandler{
		BookRepository: bookRepository,
		Validate:       validate,
	}
}

func (handler *BookHandler) Create(c *fiber.Ctx) error {
	requestBody := new(dto.BookCreateRequest)

	if err := c.BodyParser(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := handler.Validate.Struct(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve claims from token"})
	}

	userRole := claims["role"].(string)

	if userRole != "Admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "need admin role to perform this action"})
	}

	book := &entity.Book{
		Name:          requestBody.Name,
		Genre:         requestBody.Genre,
		Author:        requestBody.Author,
		PublishedDate: requestBody.PublishedDate,
		Stock:         requestBody.Stock,
		Price:         requestBody.Price,
	}

	if err := handler.BookRepository.Create(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseBody := dto.BookCreateResponse{
		ID:            book.ID,
		Name:          book.Name,
		Genre:         book.Genre,
		Author:        book.Author,
		PublishedDate: book.PublishedDate,
		Stock:         book.Stock,
		Price:         book.Price,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully added new book",
		"data":    responseBody,
	})
}

func (handler *BookHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	requestBody := new(dto.BookUpdateRequest)

	if err := c.BodyParser(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := handler.Validate.Struct(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve claims from token"})
	}

	userRole := claims["role"].(string)

	if userRole != "Admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "need admin role to perform this action"})
	}

	book, err := handler.BookRepository.FindById(bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	book.Name = requestBody.Name
	book.Genre = requestBody.Genre
	book.Author = requestBody.Author
	book.PublishedDate = requestBody.PublishedDate
	book.Stock = requestBody.Stock
	book.Price = requestBody.Price

	if err := handler.BookRepository.Update(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully update book",
		"data":    book,
	})
}

func (handler *BookHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve claims from token"})
	}

	userRole := claims["role"].(string)

	if userRole != "Admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "need admin role to perform this action"})
	}

	if err := handler.BookRepository.Delete(bookId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Successfully deleted book with ID %d", bookId)})
}

func (handler *BookHandler) FindAll(c *fiber.Ctx) error {
	books, err := handler.BookRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}

func (handler *BookHandler) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	book, err := handler.BookRepository.FindById(bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(book)
}
