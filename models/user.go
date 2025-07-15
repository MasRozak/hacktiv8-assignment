package models

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Bio      string `json:"bio" db:"bio"`
}

type Post struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type Like struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	PostID string `json:"post_id" db:"post_id"`
}

type Comment struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	PostID    string `json:"post_id" db:"post_id"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type Follow struct {
	ID          string `json:"id" db:"id"`
	FollowerID  string `json:"follower_id" db:"follower_id"`
	FollowingID string `json:"following_id" db:"following_id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error"`
}
