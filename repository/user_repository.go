package repository

import (
	"dgw-technical-test/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(user *entity.User) error
	FindUserByUsername(username string) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
}

type UserRepositoryImpl struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

func (repository *UserRepositoryImpl) Register(user *entity.User) error {
	query := "INSERT INTO Users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"

	if err := repository.DB.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryImpl) FindUserByUsername(username string) (*entity.User, error) {
	query := "SELECT * FROM Users WHERE username = $1"

	user := new(entity.User)
	if err := repository.DB.Get(user, query, username); err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindUserByEmail(email string) (*entity.User, error) {
	query := "SELECT * FROM Users WHERE email = $1"

	user := new(entity.User)
	if err := repository.DB.Get(user, query, email); err != nil {
		return nil, err
	}

	return user, nil
}
