package handler

import (
	"dgw-technical-test/dto"
	"dgw-technical-test/entity"
	"dgw-technical-test/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	return nil
}

func (handler *BookHandler) Delete(c *fiber.Ctx) error {
	return nil
}

func (handler *BookHandler) FindAll(c *fiber.Ctx) error {
	return nil
}

func (handler *BookHandler) FindById(c *fiber.Ctx) error {
	return nil
}
