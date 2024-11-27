package repository

import (
	"github.com/NekKkMirror/medods-tz.git/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save saves a new user.
func (repo *UserRepository) Save(id string, email string) error {
	_, err := repo.db.Exec("INSERT INTO users (id, email) VALUES ($1, $2)", id, email)
	return err
}

// GetById retrieves a user by their unique identifier from the database.
func (repo *UserRepository) GetById(id string) (*model.User, error) {
	var user model.User
	err := repo.db.Get(&user, "SELECT id, email FROM users WHERE id=$1", id)
	return &user, err
}
