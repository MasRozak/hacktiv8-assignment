package services

import (
	"errors"
	"strings"
	"time"

	"social-media-api/database"
	"social-media-api/models"

	"github.com/google/uuid"
)

type FollowService struct{}

func NewFollowService() *FollowService {
	return &FollowService{}
}

func (s *FollowService) CreateFollow(follow *models.Follow) error {
	if follow.FollowerID == follow.FollowingID {
		return errors.New("cannot follow yourself")
	}

	var followerExists bool
	checkFollowerQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkFollowerQuery, follow.FollowerID).Scan(&followerExists)
	if err != nil {
		return err
	}
	if !followerExists {
		return errors.New("follower user not found")
	}

	var followingExists bool
	checkFollowingQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err = database.DB.QueryRow(checkFollowingQuery, follow.FollowingID).Scan(&followingExists)
	if err != nil {
		return err
	}
	if !followingExists {
		return errors.New("following user not found")
	}

	follow.ID = uuid.New().String()
	follow.CreatedAt = time.Now().Format(time.RFC3339)

	query := `INSERT INTO follows (id, follower_id, following_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, follow.ID, follow.FollowerID, follow.FollowingID, follow.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.New("already following this user")
		}
		return err
	}

	return nil
}

func (s *FollowService) DeleteFollow(followerID, followingID string) error {
	var followExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`
	err := database.DB.QueryRow(checkQuery, followerID, followingID).Scan(&followExists)
	if err != nil {
		return err
	}
	if !followExists {
		return errors.New("follow relationship not found")
	}

	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	_, err = database.DB.Exec(query, followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func (s *FollowService) GetFollowers(userID string) ([]models.Follow, error) {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	query := `SELECT id, follower_id, following_id, created_at FROM follows WHERE following_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.Follow
	for rows.Next() {
		var follow models.Follow
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.CreatedAt)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follow)
	}

	return followers, nil
}

func (s *FollowService) GetFollowing(userID string) ([]models.Follow, error) {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := database.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	query := `SELECT id, follower_id, following_id, created_at FROM follows WHERE follower_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []models.Follow
	for rows.Next() {
		var follow models.Follow
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.CreatedAt)
		if err != nil {
			return nil, err
		}
		following = append(following, follow)
	}

	return following, nil
}
