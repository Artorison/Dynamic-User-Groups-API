package repository

import (
	"API/internal/models"
	"database/sql"
	"errors"
	"log/slog"
)

//go:generate mockery --name=UserRepository --output=mocks --outpkg=mocks
type UserRepository interface {
	GetAllUsersDB() ([]models.Users, error)
	CreateUserDB(user *models.Users) error
	DeleteUserDB(userID int64) error
	CheckUserExists(userID int64) (bool, error)
}

type UserRepositoryDB struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{DB: db}
}

func (r *UserRepositoryDB) GetAllUsersDB() ([]models.Users, error) {
	query := `SELECT * FROM users`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.Users
	for rows.Next() {
		var user models.Users
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryDB) CreateUserDB(user *models.Users) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID < 1 {
		return errors.New("invalid userID")
	}
	const op = "internal/repositories/CreateUserDB"
	query := `
        INSERT INTO users (id, name)
        VALUES ($1, $2);
    `
	if _, err := r.DB.Exec(query, user.ID, user.Name); err != nil {
		slog.String("op", op)
		return err
	}
	return nil
}

func (r *UserRepositoryDB) DeleteUserDB(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	if _, err := r.DB.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryDB) CheckUserExists(userID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	if err := r.DB.QueryRow(query, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
