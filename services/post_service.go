package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"social-media-api/database"
	"social-media-api/models"

	"github.com/google/uuid"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) CreatePost(post *models.Post) error {
	if post.Content == "" {
		return errors.New("content tidak boleh kosong")
	}

	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, post.UserID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user harus valid")
	}

	post.ID = uuid.New().String()
	post.CreatedAt = time.Now().Format(time.RFC3339)

	query := `INSERT INTO posts (id, user_id, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, post.ID, post.UserID, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	query := `SELECT id, user_id, content, created_at FROM posts ORDER BY created_at DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostService) GetPostsWithFilters(userID, keyword string) ([]models.Post, error) {
	var query string
	var args []interface{}

	baseQuery := `SELECT id, user_id, content, created_at FROM posts WHERE 1=1`

	if userID != "" {
		var userExists bool
		checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
		err := database.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
		if err != nil {
			return nil, err
		}
		if !userExists {
			return nil, errors.New("user not found")
		}

		baseQuery += ` AND user_id = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, userID)
	}

	if keyword != "" {
		baseQuery += ` AND content ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+keyword+"%")
	}

	query = baseQuery + ` ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostService) GetPostByID(id string) (*models.Post, error) {
	var post models.Post
	query := `SELECT id, user_id, content, created_at FROM posts WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostService) GetPostsByUserID(userID string) ([]models.Post, error) {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	query := `SELECT id, user_id, content, created_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostService) DeletePost(id string) error {
	var postExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err := database.DB.QueryRow(checkQuery, id).Scan(&postExists)
	if err != nil {
		return err
	}
	if !postExists {
		return errors.New("post not found")
	}

	query := `DELETE FROM posts WHERE id = $1`
	_, err = database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
