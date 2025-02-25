package repository

import (
	"dgw-technical-test/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Login(user *entity.User) error
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

func (ur *UserRepositoryImpl) Register(user *entity.User) error {
	query := "INSERT INTO Users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"

	if err := ur.DB.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepositoryImpl) Login(user *entity.User) error {
	return nil
}

func (ur *UserRepositoryImpl) FindUserByUsername(username string) (*entity.User, error) {
	query := "SELECT * FROM Users WHERE username = $1"

	user := new(entity.User)
	if err := ur.DB.Get(user, query, username); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) FindUserByEmail(email string) (*entity.User, error) {
	query := "SELECT * FROM Users WHERE email = $1"

	user := new(entity.User)
	if err := ur.DB.Get(user, query, email); err != nil {
		return nil, err
	}

	return user, nil
}
