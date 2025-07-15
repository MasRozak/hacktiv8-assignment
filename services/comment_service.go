package services

import (
	"errors"
	"time"

	"social-media-api/database"
	"social-media-api/models"

	"github.com/google/uuid"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) CreateComment(comment *models.Comment) error {
	if comment.Content == "" {
		return errors.New("content tidak boleh kosong")
	}

	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, comment.UserID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	var postExists bool
	checkPostQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err = database.DB.QueryRow(checkPostQuery, comment.PostID).Scan(&postExists)
	if err != nil {
		return err
	}
	if !postExists {
		return errors.New("post not found")
	}

	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now().Format(time.RFC3339)

	query := `INSERT INTO comments (id, user_id, post_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = database.DB.Exec(query, comment.ID, comment.UserID, comment.PostID, comment.Content, comment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentService) GetCommentsByPostID(postID string) ([]models.Comment, error) {
	var postExists bool
	checkPostQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err := database.DB.QueryRow(checkPostQuery, postID).Scan(&postExists)
	if err != nil {
		return nil, err
	}
	if !postExists {
		return nil, errors.New("post not found")
	}

	query := `SELECT id, user_id, post_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC`
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
