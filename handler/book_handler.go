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

// @Summary      Create book
// @Description  Add new book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param        request  body      dto.BookCreateRequest  true  "Create Request"
// @Success      201      {object}  dto.BookCreateResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /books [post]
// @Security     Bearer
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

// @Summary      Update book
// @Description  Update book with id
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param        request  body      dto.BookUpdateRequest  true  "Update Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /books/:id [put]
// @Security     Bearer
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

// @Summary      Delete book
// @Description  Delete book with id
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /books/:id [delete]
// @Security     Bearer
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

// @Summary      Get all books
// @Description  Retrieves a list of books
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Success      200      {array}   entity.Book
// @Failure      500      {object}  map[string]string
// @Router       /books [get]
// @Security     Bearer
func (handler *BookHandler) FindAll(c *fiber.Ctx) error {
	books, err := handler.BookRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}

// @Summary      Get book by id
// @Description  Retrieves a book by id
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Success      200      {object}  entity.Book
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /books/:id [get]
// @Security     Bearer
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
