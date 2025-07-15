package services

import (
	"database/sql"
	"errors"
	"strings"

	"social-media-api/database"
	"social-media-api/models"
	"social-media-api/utils"

	"github.com/google/uuid"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !utils.IsValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	user.ID = uuid.New().String()

	query := `INSERT INTO users (id, username, email, bio) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, user.ID, user.Username, user.Email, user.Bio)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("username or email already exists")
		}
		return err
	}

	return nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, bio FROM users ORDER BY username`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Bio)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, bio FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !utils.IsValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	var existingUser models.User
	checkQuery := `SELECT id FROM users WHERE id = $1`
	err := database.DB.QueryRow(checkQuery, id).Scan(&existingUser.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	query := `UPDATE users SET username = $1, email = $2, bio = $3 WHERE id = $4`
	_, err = database.DB.Exec(query, user.Username, user.Email, user.Bio, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("username or email already exists")
		}
		return err
	}

	user.ID = id
	return nil
}

func (s *UserService) DeleteUser(id string) error {
	var existingUser models.User
	checkQuery := `SELECT id FROM users WHERE id = $1`
	err := database.DB.QueryRow(checkQuery, id).Scan(&existingUser.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err = database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
