package services

import (
	"errors"
	"strings"

	"social-media-api/database"
	"social-media-api/models"

	"github.com/google/uuid"
)

type LikeService struct{}

func NewLikeService() *LikeService {
	return &LikeService{}
}

func (s *LikeService) CreateLike(like *models.Like) error {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, like.UserID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	var postExists bool
	checkPostQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err = database.DB.QueryRow(checkPostQuery, like.PostID).Scan(&postExists)
	if err != nil {
		return err
	}
	if !postExists {
		return errors.New("post not found")
	}

	like.ID = uuid.New().String()

	query := `INSERT INTO likes (id, user_id, post_id) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(query, like.ID, like.UserID, like.PostID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.New("satu user hanya boleh like satu post satu kali")
		}
		return err
	}

	return nil
}

func (s *LikeService) GetLikesByPostID(postID string) ([]models.Like, error) {
	var postExists bool
	checkPostQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err := database.DB.QueryRow(checkPostQuery, postID).Scan(&postExists)
	if err != nil {
		return nil, err
	}
	if !postExists {
		return nil, errors.New("post not found")
	}

	query := `SELECT id, user_id, post_id FROM likes WHERE post_id = $1`
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (s *LikeService) GetLikesByUserID(userID string) ([]models.Like, error) {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	query := `SELECT id, user_id, post_id FROM likes WHERE user_id = $1`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}
