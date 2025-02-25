package handler

import (
	"database/sql"
	"dgw-technical-test/dto"
	"dgw-technical-test/entity"
	"dgw-technical-test/repository"
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserHandler(userRepository repository.UserRepository, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

// @Summary      Register a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserRegisterRequest  true  "Register Request"
// @Success      201      {object}  dto.UserRegisterResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/register [post]
func (handler *UserHandler) Register(c *fiber.Ctx) error {
	requestBody := new(dto.UserRegisterRequest)

	if err := c.BodyParser(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := handler.Validate.Struct(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingUser, err := handler.UserRepository.FindUserByUsername(requestBody.Username)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "username or email already exists"})
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existingUser, err = handler.UserRepository.FindUserByEmail(requestBody.Email)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "username or email already exists"})
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if requestBody.Role != "Admin" && requestBody.Role != "Customer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user := &entity.User{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: string(hashPassword),
		Role:     requestBody.Role,
	}

	if err := handler.UserRepository.Register(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseBody := dto.UserRegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully register new user",
		"data":    responseBody,
	})
}

// @Summary      User login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserLoginRequest  true  "Login Request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/login [post]
func (uh *UserHandler) Login(c *fiber.Ctx) error {
	requestBody := new(dto.UserLoginRequest)

	if err := c.BodyParser(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := uh.Validate.Struct(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := uh.UserRepository.FindUserByUsername(requestBody.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully login",
		"token":   tokenString,
	})
}
