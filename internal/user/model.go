package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	Email          string    `db:"email" json:"-"`
	PasswordHash   string    `db:"password_hash" json:"-"`
	DisplayName    string    `db:"display_name" json:"display_name"`
	Bio            string    `db:"bio" json:"bio"`
	AvatarURL      string    `db:"avatar_url" json:"avatar_url"`
	FollowersCount int       `db:"followers_count" json:"followers_count"`
	FollowingCount int       `db:"following_count" json:"following_count"`
	TweetsCount    int       `db:"tweets_count" json:"tweets_count"`
	IsVerified     bool      `db:"is_verified" json:"is_verified"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserDTO struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required"`
}

type UpdateUserDTO struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio" binding:"max=160"`
	AvatarURL   string `json:"avatar_url"`
}

type Response struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"display_name"`
	Bio            string    `json:"bio"`
	AvatarURL      string    `json:"avatar_url"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	TweetsCount    int       `json:"tweets_count"`
	IsVerified     bool      `json:"is_verified"`
}
